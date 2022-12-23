import './App.css';

import axios from 'axios';
import React from 'react';
import ImageUploading from 'react-images-uploading';

import { Link } from 'react-router-dom';

import Button from '@mui/material/Button';
import Alert from '@mui/material/Alert'
import {ReactSession} from "react-client-session";
import {v4} from "uuid";

export default function Detect() {
    const [image, setImage] = React.useState([]);
    const [res, setRes] = React.useState([]);
    const maxNumber = 1;
    const onChange = (imageList) => {
        // console.log(imageList, addUpdateIndex);
        setImage(imageList);
        setRes(null);
    };

    ReactSession.setStoreType("localStorage");
    const userId = ReactSession.get("userId") ? ReactSession.get("userId") : v4() ;
    ReactSession.set("userId", userId);


    const handleSubmit = async(event) => {
        event.preventDefault();
        const formData = new FormData();
        formData.append("image", image[0].file);
        axios({
            method: "post",
            url: "https://5a54-89-110-26-237.eu.ngrok.io/detect",
            data: formData,
            headers: {
                Accept: 'application/json',
                'Content-Type': 'multipart/form-data',
            },
        }).then(resp =>
            setRes(resp.data)
        )
        .catch(error =>
            console.error(`Error: ${error}`)
        );
    }

    return (
        <div className="Detect">
            <header className="Detect-header">
                <p>
                    Различать картинки)))
                </p>
            </header>
            <div className="Detect">
                <ImageUploading value={image} onChange={onChange} maxNumber={maxNumber} dataURLKey="data_url" acceptType={["jpg"]}>
                    {({
                          imageList,
                          onImageUpload,
                          onImageRemoveAll,
                          isDragging,
                          dragProps,
                          errors
                      }) => (
                        <div className="upload__image-wrapper">
                            {imageList.map((image, index) => (
                                <div key={index} className="image-item">
                                    <img src={image.data_url} alt="" height="400"  />
                                </div>
                            ))}

                            {errors && <div>
                                {errors.maxNumber && <Alert severity="error">Можно загруть одно изображение</Alert>}
                                {errors.acceptType && <Alert severity="error">Формат файла не поддерживается</Alert>}
                                {errors.maxFileSize && <Alert severity="error">Размер файла слишком большой</Alert>}
                                {errors.resolution && <Alert severity="error">Невозможно обработать файл</Alert>}
                            </div>}&nbsp;

                            <p>{res}</p>

                            <div className="button-group" variant="contained">
                                <Button style={isDragging ? { color: "red" } : null} onClick={onImageUpload} {...dragProps}
                                        variant="contained" color="success" size="large">
                                    Место для Drag and drop
                                </Button>
                                &nbsp;
                                <Button onClick={onImageRemoveAll}
                                        variant="contained" color="success" size="large">
                                    Remove
                                </Button>
                            </div>
                            &nbsp;
                            <div className="button-group">
                                <Button onClick={handleSubmit} variant="contained">Получить тип</Button>
                            </div>
                            &nbsp;
                        </div>
                    )}
                </ImageUploading>
            </div>
            <p>user: {userId}</p>
            <div className="Footer">
                <Link to="/access">Размечать</Link>
            </div>
        </div>
    );
}