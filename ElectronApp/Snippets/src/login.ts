import 'bootstrap/dist/css/bootstrap.min.css';
import 'bootstrap';
import 'bootstrap-select';
import 'bootstrap-select/dist/css/bootstrap-select.min.css';
import axios from "axios";

const LOGIN_URL = 'http://localhost:8080/login';



var loginBtn = document.getElementById("loginBtn");
var registerBtn = document.getElementById("registerBtn");
var userNameLBL = document.getElementById("username");

console.log(loginBtn)

window.addEventListener('focus', () => userNameLBL.focus())
//window.addEventListener('DOMContentLoaded', () => initialize());

loginBtn.addEventListener("click", function() {

    
    var username = (document.getElementById("username") as HTMLInputElement).value;
    var password = (document.getElementById("password") as HTMLInputElement).value;
    authenticate(username, password);
});

function initialize() {

}

function authenticate(username: string, password: string) {
    const ipc = require('electron').ipcRenderer;
    ipc.sendSync('entry-accepted', 'hideLoginForm')

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
    
    }, (error) => {
        console.log(error);
    });


}