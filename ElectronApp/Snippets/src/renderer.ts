/**
 * This file will automatically be loaded by webpack and run in the "renderer" context.
 * To learn more about the differences between the "main" and the "renderer" context in
 * Electron, visit:
 *
 * https://electronjs.org/docs/tutorial/application-architecture#main-and-renderer-processes
 *
 */

import "./index.css";
import "bootstrap/dist/css/bootstrap.min.css";
import "bootstrap";
import "bootstrap-select";
import "bootstrap-select/dist/css/bootstrap-select.min.css";
import axios from "axios";
// eslint-disable-next-line import/no-unresolved
import * as monaco from "monaco-editor";
import * as snippet from "./dto/snippetDto";
import * as user from "./dto/userDto";
import { clipboard } from "electron";
import jwtDecode from "jwt-decode";
import $ from "jquery";

/* ********************
 * Declarations
 * ********************/


const restApiConnectionString = "http://snippets-env-1.eba-urnkpp3r.eu-central-1.elasticbeanstalk.com/";
const defaultLanguage = "markdown";
let currentUser: user.UserDto = null;
let currentSnippet: snippet.SnippetDto = null;
let currentSnippetModified = false;
let jwtAuthToken = "";
const jwtHeaderConfig = {
  headers: {
    Authorization: "" + jwtAuthToken,
  },
};

const languages = [
  "abap",
  "apex",
  "azcli",
  "bat",
  "cameligo",
  "clojure",
  "coffee",
  "cpp",
  "csharp",
  "csp",
  "css",
  "dockerfile",
  "fsharp",
  "go",
  "graphql",
  "handlebars",
  "html",
  "ini",
  "java",
  "javascript",
  "json",
  "kotlin",
  "less",
  "lua",
  "markdown",
  "mips",
  "msdax",
  "mysql",
  "objective-c",
  "pascal",
  "pascaligo",
  "perl",
  "pgsql",
  "php",
  "postiats",
  "powerquery",
  "powershell",
  "pug",
  "python",
  "r",
  "razor",
  "redis",
  "redshift",
  "restructuredtext",
  "ruby",
  "rust",
  "sb",
  "scheme",
  "scss",
  "shell",
  "solidity",
  "sophia",
  "sql",
  "st",
  "swift",
  "tcl",
  "twig",
  "typescript",
  "vb",
  "xml",
  "yaml",
];

const editor = monaco.editor.create(
  document.getElementById("monaco-container"),
  {
    theme: "vs-dark",
    scrollBeyondLastLine: false,
  }
);

const model = monaco.editor.createModel("", "markdown");

const selector = $("#languageSelector").get(0) as HTMLSelectElement;

/* ********************
 * Functions
 * ********************/

function logout(): void {
  const emptyText = "";

  jwtAuthToken = "";
  currentUser = null;
  currentSnippet = null;
  currentSnippetModified = false;

  $("#loginAlert").hide();
  $("#monacoSnippetName").text(emptyText);
  $("#monacoSaveHint").text(emptyText);
  model.setValue(emptyText);
  selector.selectedIndex = languages.indexOf(defaultLanguage);
  $("#userNameLink").text(emptyText);
  $("#snippetList").empty();

  ($("#loginModal") as any).modal("show");
}

function updateDimensions(): void {
  editor.layout();
}

function deriveUsernameFromEmail(usermail: string): string {
  const splittedMail = usermail.split("@");
  return splittedMail[0];
}

function saveCurrentSnippetFromModel(): void {
  currentSnippet.code = model.getValue();
  currentSnippet.title = $("#monacoSnippetName").first().text().toString();
  currentSnippet.language = languages[selector.selectedIndex];
}

function saveToClipboard(): void {
  if (currentSnippet != null) {
    clipboard.writeText(`${JSON.stringify(currentSnippet)}`);
  }
}

function setSnippetSaveNotifer(text: string): void {
  $("#monacoSaveHint").text(text);
}

/**
 * @param {String} HTML representing a single element
 * @return {Element}
 */
function htmlToElement(html: string): ChildNode {
  const template = document.createElement("template");
  html = html.trim(); // Never return a text node of whitespace as the result
  template.innerHTML = html;
  return template.content.firstChild;
}

function loadLanguages(): void {
  let count = 0;

  for (const language of languages) {
    const opt = document.createElement("option");
    opt.value = count.toString();
    opt.text = language;

    selector.appendChild(opt);
    count += 1;
  }

  selector.addEventListener("change", function () {
    const language = languages[parseInt(selector.value)];
    monaco.editor.setModelLanguage(model, language);
    currentSnippet.language = language;
    setSnippetSaveNotifer("*not saved");
    currentSnippetModified = true;
  });

  // set default language
  selector.selectedIndex = languages.indexOf(defaultLanguage);
}

function setModelWithLanguage(loadedSnippet: snippet.SnippetDto): void {
  // set current snippet
  currentSnippet = loadedSnippet;

  model.setValue(loadedSnippet.code);
  monaco.editor.setModelLanguage(model, loadedSnippet.language);
  selector.selectedIndex = languages.indexOf(loadedSnippet.language);
  $(".selectpicker").selectpicker("refresh");
  $("#monacoSnippetName").text(loadedSnippet.title);
}

async function loadUserAsync(
  connectionString: string,
  usermail: string,
  password: string
): Promise<user.UserDto> {
  const username = deriveUsernameFromEmail(usermail);

  return axios
    .post(`${connectionString}`, {
      user_id: "",
      username: username,
      mail: usermail,
      password: password,
    })
    .then((response) => {
      // TODO: check if its really a jwt token
      jwtAuthToken = response.data;
      jwtHeaderConfig.headers.Authorization = "" + jwtAuthToken;

      const jwtData = jwtDecode(jwtAuthToken);

      if (!user.isUserDto(jwtData)) {
        console.error(
          `Invalid request - expected UserDTO got:\n${response.data}`
        );
        console.info(response);
        return null;
      } else {
        return jwtData;
      }
    })
    .catch((error) => {
      console.log(`Error while loading user ${usermail}`, error);
      return null;
    });
}

async function loadSnippetAsync(
  connectionString: string,
  userId: string,
  snippetId: string
): Promise<snippet.SnippetDto> {
  return axios
    .get(
      `${connectionString}user/${userId}/snippets/${snippetId}`,
      jwtHeaderConfig
    )
    .then((response) => {
      if (!snippet.isSnippetDto(response.data)) {
        console.error(`Invalid request - expected SnippetDTO`);
        console.info(response);
        return null;
      } else {
        return response.data;
      }
    })
    .catch((error) => {
      console.log(
        `Error while loading snippet with id ${snippetId} for user ${userId}`,
        error
      );
      return null;
    });
}

async function loadSnippetsAsync(
  connectionString: string,
  userId: string
): Promise<snippet.SnippetDto[]> {
  return axios
    .get(`${connectionString}user/${userId}/snippets`, jwtHeaderConfig)
    .then((response) => {
      if (response.data.length == 0) {
        return [] as snippet.SnippetDto[];
      }

      if (!snippet.isSnippetDto(response.data[0])) {
        console.error(`Invalid request - expected SnippetDTO`);
        return [] as snippet.SnippetDto[];
      } else {
        const snippets: snippet.SnippetDto[] = Object.keys(response.data).map(
          (i) => response.data[i]
        );
        return snippets;
      }
    })
    .catch((error) => {
      console.log(`Error while loading snippets for user ${userId}`, error);
      return [] as snippet.SnippetDto[];
    });
}

async function createNewSnippetAsync(
  connectionString: string,
  userId: string
): Promise<snippet.SnippetDto> {
  const data = {
    owner: currentUser.user_id,
    title: "New Snippet",
    category: "",
    code: "",
    language: "markdown",
  };

  return axios
    .post(`${connectionString}user/${userId}/snippets`, data, {
      headers: jwtHeaderConfig.headers,
    })
    .then((response) => {
      if (!snippet.isSnippetDto(response.data)) {
        console.error(`Invalid request - expected SnippetDTO`);
        console.info(response);
        return null;
      } else {
        return response.data;
      }
    })
    .catch((error) => {
      console.log(`Error while creating snippets for user ${userId}`, error);
      return null;
    });
}

async function updateSnippetAsync(
  connectionString: string,
  userId: string
): Promise<snippet.SnippetDto> {
  return axios
    .put(
      `${connectionString}user/${userId}/snippets/${currentSnippet.snippet_id}`,
      JSON.stringify(currentSnippet),
      { headers: jwtHeaderConfig.headers }
    )
    .then((response) => {
      if (!snippet.isSnippetDto(response.data)) {
        console.error(`Invalid request - expected SnippetDTO`);
        console.info(response);
        return null;
      } else {
        currentSnippetModified = false;
        return response.data;
      }
    })
    .catch((error) => {
      console.log(
        `Error while updating snippet for user ${userId} snippet: ${currentSnippet.snippet_id}`,
        error
      );
      return null;
    });
}

async function registerSnippetLink(
  connectionString: string,
  s: snippet.SnippetDto
): Promise<void> {
  if (currentSnippet != null && currentSnippetModified) {
    await updateSnippetAsync(restApiConnectionString, currentUser.user_id);
  }

  loadSnippetAsync(connectionString, currentUser.user_id, s.snippet_id).then(
    (response) => {
      if (response != null) {
        setModelWithLanguage(response);
      }
    }
  );
}

function createSnippetLinks(
  connectionString: string,
  snippets: snippet.SnippetDto[]
): void {
  const ul = $("#snippetList").first();

  for (const s of snippets) {
    const html = `<li class="list-group-item" style="padding: 0em;"><a id="${s.snippet_id}" class="nav-link" href="#">${s.title}</a></li>`;

    const li = htmlToElement(html);
    li.addEventListener("click", async () => {
      registerSnippetLink(connectionString, s);
    });

    ul.append(li as HTMLElement);
  }
}

function loginAndRegisterResponseHandler(response: user.UserDto): void {
  if (response != null) {
    currentUser = response;
    ($("#loginModal") as any).modal("hide").data("#loginModal", null);
    $("#userNameLink").text("Hi " + currentUser.username);

    loadSnippetsAsync(restApiConnectionString, currentUser.user_id).then(
      (snippetResponse3) => {
        if (snippetResponse3.length > 0) {
          createSnippetLinks(restApiConnectionString, snippetResponse3);
        }
      }
    );
  } else {
    const usermail = $("#usermail");
    const password = $("#password");

    usermail.addClass("is-invalid");
    password.addClass("is-invalid");

    $("#loginAlert").text(
      "Either invalid password or no user with that email found."
    );
    $("#loginAlert").show();
  }
}

function loadMainApplication(path: string): void {
  const requestUrl = restApiConnectionString + path;
  const usermail = $("#usermail");
  const password = $("#password");

  let invalid = false;
  $("#loginAlert").hide();

  if ((usermail.val() as string) === "") {
    usermail.addClass("is-invalid");
    invalid = true;
  } else {
    usermail.removeClass("is-invalid");
  }

  if ((password.val() as string) === "") {
    password.addClass("is-invalid");
    invalid = true;
  } else {
    password.removeClass("is-invalid");
  }

  if (!invalid) {
    loadUserAsync(
      requestUrl,
      usermail.val() as string,
      password.val() as string
    ).then(loginAndRegisterResponseHandler);
  }
}

function addLoginListener(): void {
  $("#loginAlert").hide();

  // Handle login
  $("#loginBtn").click(() => {
    loadMainApplication("login");
  });

  // Handle Register
  $("#registerBtn").click(() => loadMainApplication("user"));
}

function addLogoutListener(): void {
  $("#logoutBtn").click(() => logout());
}

function initializeMonacoEditor(): void {
  editor.setModel(model);

  // Add bindings
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  editor.addCommand(monaco.KeyMod.CtrlCmd | monaco.KeyCode.KEY_S, function () {
    saveCurrentSnippetFromModel();
    updateSnippetAsync(restApiConnectionString, currentUser.user_id)
      .then((response) => {
        setSnippetSaveNotifer("");
      })
      .catch((error) => {
        console.log(
          `Error while saving snippet for user ${currentUser.user_id}, snippet: ${currentSnippet.snippet_id}`,
          error
        );
        setSnippetSaveNotifer("*not saved");
      });
  });

  editor.onDidChangeModelContent((event) => {
    if (currentSnippet != null && currentSnippet.code != model.getValue()) {
      setSnippetSaveNotifer("*not saved");
      currentSnippetModified = true;
      currentSnippet.code = model.getValue();
    } else {
      setSnippetSaveNotifer("");
    }
  });

  editor.layout();
}

function intializeMainApplication(): void {
  addLoginListener();
  addLogoutListener();
  initializeMonacoEditor();
  loadLanguages();

  ($("[data-toggle=popover]") as any).popover();
  $("#monacoSnippetName").on("click", () => {
    if (currentSnippet != null) {
      $("#snippetUpdateName").val(currentSnippet.title);
      ($("#changeNameModal") as any).modal("show");
    }
  });


  $("#shareButton").on("click", () => {
    saveToClipboard();
  });

  $("#updateSnippetButton").on("click", () => {
    const newSnippetName = $("#snippetUpdateName").val() as string;

    $("#monacoSnippetName").text(newSnippetName);
    $(`#${currentSnippet.snippet_id}`).text(newSnippetName);
    ($("#changeNameModal") as any).modal("hide");

    currentSnippet.title = newSnippetName;
    updateSnippetAsync(restApiConnectionString, currentUser.user_id);
  });

  $("#snippetCreationLink").on("click", () => {
    createNewSnippetAsync(restApiConnectionString, currentUser.user_id).then(
      (snippetResponse) => {
        const ul = document.getElementById("snippetList");
        const html = `<li class="list-group-item" style="padding: 0em;"><a id="${snippetResponse.snippet_id}" class="nav-link" href="#">${snippetResponse.title}</a></li>`;
        const li = htmlToElement(html);

        li.addEventListener("click", async () => {
          if (currentSnippet != null && currentSnippetModified) {
            await updateSnippetAsync(
              restApiConnectionString,
              currentUser.user_id
            );
          }

          loadSnippetAsync(
            restApiConnectionString,
            currentUser.user_id,
            snippetResponse.snippet_id
          ).then((snippetResponse2) => {
            if (snippetResponse2 != null) {
              setModelWithLanguage(snippetResponse2);
            }
          });
        });

        ul.appendChild(li);
      }
    );
  });
}

// Load Login Modal on start up
$(window).on("load", function () {
  ($("#loginModal") as any).modal("show");
});

window.addEventListener("resize", updateDimensions.bind(this));
window.onclose = function (): void {
  window.removeEventListener("resize", updateDimensions.bind(this));
};

intializeMainApplication();
