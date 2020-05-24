import 'bootstrap/dist/css/bootstrap.min.css';
import 'bootstrap';
import 'bootstrap-select';
import 'bootstrap-select/dist/css/bootstrap-select.min.css';
import axios from "axios";

const LOGIN_URL = 'http://localhost:8080/login';


//var loginBtn = document.getElementById("snippetLoginBtn");

/*loginBtn.addEventListener("click", function() {

    let $ = require('jquery')
    const ipc = require('electron').ipcRenderer;

    ipc.sendSync('entry-accepted', 'ping')
    //var username = (document.getElementById("inputUserName") as HTMLInputElement).value;
    //var password = (document.getElementById("inputPassword") as HTMLInputElement).value;
    //authenticate(username, password);
});


function authenticate(username: string, password: string) {
    console.log(username + " . " + password);
    // TODO: hash password
    
    axios.post('http://127.0.0.1:8080/login', 
        {
            username : username,
            password : password,
        

    }).then(response => {
        console.log(response);
        //TODO: handle response (JWT)
        //TODO mb use a logindialog
        //document.location.href = './index.html';
    
    }, (error) => {
        console.log(error);
        //document.location.href = 'index.html';
    });

    */
}