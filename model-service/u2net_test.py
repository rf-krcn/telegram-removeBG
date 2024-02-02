import numpy as np
import os
from skimage import io, transform
import torch
import torchvision
from torch.autograd import Variable
import torch.nn.functional as F
from torch.utils.data import Dataset, DataLoader
from torchvision import transforms
from PIL import Image
import glob
from flask import Flask, request, jsonify, send_file
from io import BytesIO
from waitress import serve

from data_loader import RescaleT
from data_loader import ToTensor
from data_loader import ToTensorLab
from data_loader import SalObjDataset
from model import U2NET  # full size version 173.6 MB
from model import U2NETP  # small version u2net 4.7 MB

model_dir = os.path.join(os.getcwd(), 'saved_models', 'u2net','u2net' + '.pth')
prediction_dir = os.path.join(os.getcwd(), 'test_data', 'u2net' + '_results' + os.sep)
net = U2NET(3,1)

if torch.cuda.is_available():
    net.load_state_dict(torch.load(model_dir))
    net.cuda()
else:
    net.load_state_dict(torch.load(model_dir, map_location='cpu'))
net.eval()

# normalize the predicted SOD probability map
def normPRED(d):
    ma = torch.max(d)
    mi = torch.min(d)

    dn = (d-mi)/(ma-mi)

    return dn

def save_output(image_name,pred,d_dir):

    predict = pred
    predict = predict.squeeze()
    predict_np = predict.cpu().data.numpy()

    im = Image.fromarray(predict_np*255).convert('RGB')
    img_name = image_name.split(os.sep)[-1]
    image = io.imread(image_name)
    imo = im.resize((image.shape[1],image.shape[0]),resample=Image.BILINEAR)

    pb_np = np.array(imo)

    aaa = img_name.split(".")
    bbb = aaa[0:-1]
    imidx = bbb[0]
    for i in range(1,len(bbb)):
        imidx = imidx + "." + bbb[i]

    imo.save(d_dir+imidx+'.png')

app = Flask(__name__)

@app.route('/process_image', methods=['POST'])
def process_image():
    try:
        # Get the image data from the request
        image_data = request.data

        # Save the image to the test_images directory
        image_dir = os.path.join(os.getcwd(), 'test_data', 'test_images')
        image_path = [os.path.join(image_dir, 'raw.JPG')]
        with open(image_path[0], 'wb') as image_file:
            image_file.write(image_data)



        test_salobj_dataset = SalObjDataset(img_name_list = image_path,
                                        lbl_name_list = [],
                                        transform=transforms.Compose([RescaleT(320),
                                                                      ToTensorLab(flag=0)])
                                        )
        test_salobj_dataloader = DataLoader(test_salobj_dataset,
                                        batch_size=1,
                                        shuffle=False,
                                        num_workers=1)



        # Process the image
        inputs_test = next(iter(test_salobj_dataloader))['image']
        inputs_test = inputs_test.type(torch.FloatTensor)

        if torch.cuda.is_available():
            inputs_test = Variable(inputs_test.cuda())
        else:
            inputs_test = Variable(inputs_test)

        d1, d2, d3, d4, d5, d6, d7 = net(inputs_test)

        # normalization
        pred = d1[:, 0, :, :]
        pred = normPRED(pred)

        # Save the results
        save_output(image_path[0], pred, prediction_dir)

        # Read the saved output image
        output_image_path = os.path.join(prediction_dir, 'raw.png')
        output_image = Image.open(output_image_path).convert("L")
        
        original = Image.open(image_path[0])
        original.putalpha(output_image)
        original.save("result_image.png")
        output_image = Image.open("result_image.png")
        
        # Convert the image to bytes
        output_image_bytes = BytesIO()
        output_image.save(output_image_bytes, format='PNG')
        output_image_bytes.seek(0)

        # Delete input and output files
        os.remove(image_path[0])
        os.remove(output_image_path)
        os.remove("result_image.png")
        # Return the processed image as bytes
        return send_file(output_image_bytes, mimetype='image/png')

    except Exception as e:
        # If an error occurs, return an error response
        error_data = {'status': 'error', 'message': str(e)}
        print(e)
        return jsonify(e), 500

if __name__ == "__main__":
    serve(app, host="0.0.0.0", port=8080)