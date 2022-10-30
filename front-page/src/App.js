import './App.css';
import React, { useEffect } from 'react';
import ImageUploading from 'react-images-uploading';
import { ReactSession } from 'react-client-session';
import { v4 as uuidv4 } from 'uuid';
import Button from '@mui/material/Button';
import Alert from '@mui/material/Alert';
import axios from 'axios';



export default function App() {
  ReactSession.setStoreType("localStorage");
  const userId = ReactSession.get("userId") ? ReactSession.get("userId") : uuidv4() ;
  ReactSession.set("userId", userId);

  const [url, setUrl] = React.useState([]);
  const [imgId, setImgId] = React.useState([]); 
  const [smallUrl, setSmallUrl] = React.useState([]);


  useEffect(() => {
    console.log("Welcome");
    getHoneyImage();
  }, [""]);

  const maxNumber = 1;
  
  const handleSubmit = (event) => {
    const jsonObj = JSON.stringify({ session_id: userId, image_id: imgId, image_url: smallUrl, type: event.target.name});
    axios({
      method: "post",
      url: "http://localhost:10000/photo",
      data: jsonObj,
      headers: {'Content-Type': 'application/json'}
    })
    .catch(error => 
      console.error(`Error: ${error}`)
    );
    getImage();
  }

  const getImage = () => {
    axios.get("http://localhost:10000/photo")
    .then((response)=>{
      setUrl(response.data.urls.regular);
      setSmallUrl(response.data.urls.thumb);
      setImgId(response.data.id);
    })
    .catch(error => 
      console.error(`Error: ${error}`)
    );
  }

  const getHoneyImage = () => {
    axios.get("http://localhost:10000/photo/honey")
    .then((response)=>{
      setUrl(response.data.urls.regular);
      setSmallUrl(response.data.urls.thumb);
      setImgId(response.data.id);
    })
    .catch(error => 
      console.error(`Error: ${error}`)
    );
  }
  
  return (
    <div className="App">
      <header className="App-header">
        <p>
          Картинки)))
        </p>
      </header>
        <div className="App">
            <div className="image-wrapper">
              <img src={url} width="800" height="400"></img>
            </div>
            &nbsp;
            <div className="button-group">
              <Button onClick={getImage} variant="contained">Получить новое изображение</Button>
            </div>
            &nbsp;
            <div className="button-group">
              <Button onClick={handleSubmit} variant="contained" color="success" size="large" name="SUNNY">Sunny</Button>
              &nbsp;
              <Button onClick={handleSubmit} variant="contained" color="success" size="large" name="CLOUDY">Cloudy</Button>
              &nbsp;
              <Button onClick={handleSubmit} variant="contained" color="success" size="large" name="RAINY">Rainy</Button>
              &nbsp;
              <Button onClick={handleSubmit} variant="contained" color="success" size="large" name="SUNRISE">Sunrise/Sunset</Button>
              &nbsp;
              <Button onClick={handleSubmit} variant="contained" color="success" size="large" name="WRONG">Not related</Button>
            </div>
        </div>
        <p>user: {userId}</p>
    </div>
  );
}

