# Electron App Steps

#### Needed software
https://nodejs.org/en/download/  

Visual Studio Extensions:
> code --install-extension dzannotti.vscode-babel-coloring  
> code --install-extension EditorConfig.EditorConfig  
> code --install-extension dbaeumer.vscode-eslint  
> code --install-extension esbenp.prettier-vscode  

#### Basic setup Electron Forge
https://www.electronforge.io/templates/typescript-+-webpack-template

Create the app:

> npx create-electron-app app-name --template=typescript-webpack

Start:

> npm start

#### Babel

Installation:  
> npm install --save-dev babel-loader @babel/core  
> npm install --save-dev @babel/preset-typescript
> npm install jquery --save-dev
> npm install -D @types/bootstrap
> npm install @types/jquery

In `webpack.rules.js` add another rule:
```typescript
  {
    test: /\.m?js$/,
    exclude: /(node_modules|bower_components)/,
    use: {
      loader: 'babel-loader',
      options: {
        presets: ['@babel/preset-typescript']
      }
    }
  },
```

#### Axios
https://www.npmjs.com/package/axios

Install:

> npm install axios

Small usage example:
```javascript
import axios from "axios";

axios
  .post("/login", {
    firstName: "Finn",
    lastName: "Williams",
  })
  .then(
    (response) => {
      console.log(response);
    },
    (error) => {
      console.log(error);
    }
  );
```
#### Add monaco editor
https://microsoft.github.io/monaco-editor/  
https://www.npmjs.com/package/monaco-editor  
https://github.com/microsoft/monaco-editor/blob/master/docs/integrate-esm.md  
https://github.com/Microsoft/monaco-editor-webpack-plugin  
https://www.npmjs.com/package/file-loader

Explanation:
We are using monaco webpack loader plugin. This allows for options to be passed into the plugin in order to select only a subset of editor features or editor languages.  Options can be passed in to MonacoWebpackPlugin. They can be used to generate a smaller editor bundle by selecting only certain languages or only certain editor features.  

Default languages are:   
`['abap', 'apex', 'azcli', 'bat', 'cameligo', 'clojure', 'coffee', 'cpp', 'csharp', 'csp', 'css', 'dockerfile', 'fsharp', 'go', 'graphql', 'handlebars', 'html', 'ini', 'java', 'javascript', 'json', 'kotlin', 'less', 'lua', 'markdown', 'mips', 'msdax', 'mysql', 'objective-c', 'pascal', 'pascaligo', 'perl', 'pgsql', 'php', 'postiats', 'powerquery', 'powershell', 'pug', 'python', 'r', 'razor', 'redis', 'redshift', 'restructuredtext', 'ruby', 'rust', 'sb', 'scheme', 'scss', 'shell', 'solidity', 'sophia', 'sql', 'st', 'swift', 'tcl', 'twig', 'typescript', 'vb', 'xml', 'yaml']`  

Default features are:  
`['accessibilityHelp', 'bracketMatching', 'caretOperations', 'clipboard', 'codeAction', 'codelens', 'colorDetector', 'comment', 'contextmenu', 'coreCommands', 'cursorUndo', 'dnd', 'find', 'folding', 'fontZoom', 'format', 'gotoError', 'gotoLine', 'gotoSymbol', 'hover', 'iPadShowKeyboard', 'inPlaceReplace', 'inspectTokens', 'linesOperations', 'links', 'multicursor', 'parameterHints', 'quickCommand', 'quickOutline', 'referenceSearch', 'rename', 'smartSelect', 'snippets', 'suggest', 'toggleHighContrast', 'toggleTabFocusMode', 'transpose', 'wordHighlighter', 'wordOperations', 'wordPartOperations']`
`  

 It is also possible to exclude certain default features prefixing them with an exclamation mark '!'.  


Current used version matrix:  
| `monaco-editor-webpack-plugin` | `monaco-editor` |
|---|---|
| `1.9.x` | `0.20.x` |


Installation:
> npm install monaco-editor  
> npm install monaco-editor-webpack-plugin
> npm install file-loader --save-dev
> npm install css-loader --save-dev 

In `webpack.plugin.js` add the plugin:  

```typescript
const MonacoWebpackPlugin = require('monaco-editor-webpack-plugin');

module.exports = [
  new MonacoWebpackPlugin(),
];
```

In `webpack.renderer.js` add `.ttf` extension:  
```typescript
module.exports = {
  module: {
    rules,
  },
  plugins: plugins,
  resolve: {
    extensions: ['.js', '.ts', '.jsx', '.tsx', '.css', '.ttf']
  },
};
```

In `webpack.rules.js` add a `file-loader` for the `.ttf` files. `file-loader` has to be installed:
```typescript
{
    test: /\.ttf$/,
    use: 'file-loader',
},
```

Usage sample:
```typescript
import * as monaco from "monaco-editor";

const model = monaco.editor.createModel(
  ["function x() {", '\tconsole.log("Hello world!");', "}"].join("\n"),
  "javascript"
);

const editor = monaco.editor.create(document.getElementById("container"), {
  theme: "vs-dark",
});

editor.setModel(model);

const myBinding = editor.addCommand(
  monaco.KeyMod.CtrlCmd | monaco.KeyCode.F10,
  function () {
    alert("CTRL + F10 pressed!");
  }
);
```
