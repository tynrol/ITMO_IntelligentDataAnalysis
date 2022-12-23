import uvicorn
if __name__ == '__main__':
    uvicorn.run('py_inference:app', host='127.0.0.1', port=33334, log_level="info")
    exit()

import torch
from torchvision import transforms
from fastapi import FastAPI, File,Form, HTTPException, Response, status
import numpy as np
from PIL import Image
import io

app = FastAPI()
model = torch.jit.load("./model.pt")
model.eval()
model.to("cpu")

_transform=transforms.Compose([
                       transforms.Resize((224,224)),
                       transforms.ToTensor(),
                       transforms.Normalize([0.5, 0.5, 0.5], [0.5, 0.5, 0.5])])


def read_img_buffer(image_data):
    img = Image.open(io.BytesIO(image_data))
    # img=img.convert('L').convert('RGB') #GREYSCALE
    if img.mode != 'RGB':
        img = img.convert('RGB')
    return img

def get_res(image_buffer):
    cls = ["cloudy","rain","shine","sunrise"]
    image = read_img_buffer(image_buffer)
    image = _transform(image).unsqueeze(0).to("cpu")
    with torch.no_grad():
        res = model(image).cpu().numpy()
    return cls[np.argmax(res)]


@app.get("/")
async def read_root():
    return {"Hello": "World"}

@app.post("/get_class")
async def global_features_get_similar_images_by_image_buffer_handler(image: bytes = File(...)):
    res  = get_res(image)
    return res