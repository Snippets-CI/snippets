/**
 * This file will automatically be loaded by webpack and run in the "renderer" context.
 * To learn more about the differences between the "main" and the "renderer" context in
 * Electron, visit:
 *
 * https://electronjs.org/docs/tutorial/application-architecture#main-and-renderer-processes
 *
 * By default, Node.js integration in this file is disabled. When enabling Node.js integration
 * in a renderer process, please be aware of potential security implications. You can read
 * more about security risks here:
 *
 * https://electronjs.org/docs/tutorial/security
 *
 * To enable Node.js integration in this file, open up `main.js` and enable the `nodeIntegration`
 * flag:
 *
 * ```
 *  // Create the browser window.
 *  mainWindow = new BrowserWindow({
 *    width: 800,
 *    height: 600,
 *    webPreferences: {
 *      nodeIntegration: true
 *    }
 *  });
 * ```
 */

import "./index.css";
import axios from "axios";
import * as monaco from "monaco-editor";

axios.get("https://api.github.com/users/mapbox").then((response) => {
  console.log(response.data);
  console.log(response.status);
  console.log(response.statusText);
  console.log(response.headers);
  console.log(response.config);
});

let model = monaco.editor.createModel(
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

console.log(
  '👋 This message is being logged by "renderer.js", included via webpack'
);
