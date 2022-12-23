import './App.css';

import axios from 'axios';
import React from 'react';

import { useEffect } from 'react';
import { ReactSession } from 'react-client-session';
import { Link } from 'react-router-dom';
import { v4 } from 'uuid';

import Button from '@mui/material/Button';

export default function Access() {
    const [url, setUrl] = React.useState([]);
    const [imgId, setImgId] = React.useState([]);
    const [smallUrl, setSmallUrl] = React.useState([]);

    ReactSession.setStoreType("localStorage");
    const userId = ReactSession.get("userId") ? ReactSession.get("userId") : v4() ;
    ReactSession.set("userId", userId);

    useEffect(() => {
      console.log("Welcome");
      getImage();
    }, [""]);

    const handleSubmit = (event) => {
        const jsonObj = JSON.stringify({ session_id: userId, image_id: imgId, image_url: smallUrl, type: event.target.name});
        axios({
            method: "post",
            url: "https://a6e5-89-110-26-237.eu.ngrok.io/photo",
            data: jsonObj,
            headers: {'Content-Type': 'application/json'}
        })
            .catch(error =>
                console.error(`Error: ${error}`)
            );
        getImage();
    }

    const getImage = () => {
        axios.get("https://a6e5-89-110-26-237.eu.ngrok.io/photo")
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
        <div className="Access">
            <header className="Access-header">
                <p>
                    Размечать картинки)))
                </p>
            </header>
            <div className="Access">
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
            <div className="Footer">
                <Link to="/">Различать</Link>
            </div>
        </div>
    );
}