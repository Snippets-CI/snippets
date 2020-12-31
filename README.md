# Snippets [![Build Status](https://travis-ci.com/AndreasRoither/Snippets.svg?branch=master)](https://travis-ci.com/AndreasRoither/Snippets)[![Bugs](https://sonarcloud.io/api/project_badges/measure?project=AndreasRoither_Snippets&metric=bugs)](https://sonarcloud.io/dashboard?id=AndreasRoither_Snippets)[![Vulnerabilities](https://sonarcloud.io/api/project_badges/measure?project=AndreasRoither_Snippets&metric=vulnerabilities)](https://sonarcloud.io/dashboard?id=AndreasRoither_Snippets)

### Motivation

Sometimes when writing code you forget what the optimal solution for a particular task like creating a button listener is. Furthermore, a lot of code has to be remembered which can be quite challenging.

### Goal

The goal of this project is to create a service where users can save code snippets. Users can log in and view their snippets that they previously added. These snippets can be easily copied and tweaked to the users liking. 

### Realization

With the Electron framework it's possible to create a cross-platform-desktop GUI application. CockroachDB is scalable and allows us to use multiple nodes. Golang will be used as an authentication service to access the database. Using docker all instances can be deployed in a fast and efficient way.

# Project parts

### Electron app
The app has been built with electron and typescript. As a packaging solution we have used [electron-forge](https://www.electronforge.io/) with [webpack](https://webpack.js.org/). Webpack allows us to bundle modules dependencies to static assets. 
We didn't use the bundeling option for modules but rather to include and pack the main application. The electron app connects to the heroku hosted rest api which will be explained later on. 

Demo:  
![](https://i.imgur.com/vg38CJI.png)  
![](https://i.imgur.com/0ZHMFGS.png)  

Snippets can be added by pressing the + sign. To change the snippet name click on name above the editor. Changing the language can be achieved by pressing the language button below the code editor.


### RestAPI
The RestAPI is created using golang and jwt. The database that the api connects to is postgres. To ensure that the rest api is working correctly tests have been created that test every available functionality for any breaking changes. Since we use travis and heroku for deployment the connection string attributes are collected using environment variables. For CI/CD configuration through environment variables is the prefered way anyway.

Test run:  
<p align="left">
  <img src="https://i.imgur.com/9uHRlNn.png" width="400"/>
</p>


### Docker & docker compose
Docker compose is used to start up three containers:  the rest api, database and pgAdmin to monitor and adjust the database.

### Heroku

Heroku is used for the rest api deployment. Together with a database provided by heroku, new app releases are using the deployed rest api for production.

<p align="left">
  <img src="https://i.imgur.com/vhXgSVl.png" width="400"/>
</p>


### Travis

Travis is used to automatically test and deploy. 
There are six stages, but not all of them are executed every time:
- sonarcloud
- test rest api
- build docker
- build docker and deploy to heroku
- make electron
- publish electron
  
The deployment to heroku and the publish electron stage are only run if a commit on master has been tagged after the commit has been pushed.

Normal push on master:  

<p align="left">
  <img src="https://i.imgur.com/eKoz9di.png" width="500"/>
</p>

Tag push on master:  
<p align="left">
  <img src="https://i.imgur.com/mpYcGyF.png" width="500"/>
</p>


### Github release

A github release is automatically create when a commit on master has been tagged:
<p align="left">
  <img src="https://i.imgur.com/FmLdJFZ.png" width="400"/>
</p>


### Tools Used
- [Electron](https://www.electronjs.org/) as GUI framework  
- [Golang](https://golang.org/) Rest API with [JWT](https://jwt.io/)  
- [Postgres](https://www.postgresql.org/) as database  
- [Docker & Docker Compose](https://www.docker.com/) for virtualization  

#### Development and extensions
- [Visual Studio Code](https://code.visualstudio.com/) for development  
    - [ESLint](https://marketplace.visualstudio.com/items?itemName=dbaeumer.vscode-eslint)  
    - [npm Inellisense](https://marketplace.visualstudio.com/items?itemName=christian-kohler.npm-intellisense)  
    - [Firefox debugger](https://marketplace.visualstudio.com/items?itemName=firefox-devtools.vscode-firefox-debug)  
    - [Prettier](https://marketplace.visualstudio.com/items?itemName=esbenp.prettier-vscode)  
