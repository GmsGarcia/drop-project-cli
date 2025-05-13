<h1 align="center">
  <a href="https://github.com/GmsGarcia/drop-project-cli"><img src="https://raw.githubusercontent.com/GmsGarcia/drop-project-cli/master/media/logo.png" alt="DropProjectCLI Logo" width="200"></a>
  <br>
  Drop Project CLI
  <br>
</h1>

<h3 align="center">Simple CLI app to upload and access Drop Project assignments</h3>

<p align="center">
  <a href="#description">Description</a>•
  <a href="#features">Features</a>•
  <a href="#installation">Installation</a>•
  <a href="#usage">Usage</a>•
  <a href="#configuration">Configuration</a>
</p>

<h2 id="description">📟 What is Drop Project and Drop Project CLI?</h2>

<a href="https://dropproject.org/">Drop Project</a> is an open-source automated assessment tool that checks student programming projects for correctness and quality.

Drop Project CLI is an cross-platform CLI app that bridges the Drop Project API with your terminal.
It allows you to review and submit your project directly from the terminal, check out the next section for all the features.

<h2 id="features">📱 Features</h2>

- Review last submission report
- Create new submissions
- ...

<h2 id="installation">📦 Installation</h2>

Download the latest release from the <a href="https://github.com/GmsGarcia/drop-project-cli/releases">releases page</a>.

<h3>Windows and MacOS</h3>

Run the executable file.

<h3>Linux</h3>

If you **ARE** running on a headless system, like WSL, you should run `sudo ./install-headless.sh`.

If you are **NOT** running on a headless system, you should run `sudo ./install.sh`.

<h2 id="usage">📝 Usage</h2>

```
  $dpc help
  
  Usage: dpc [command] [arguments]
  
  Commands:
    get|g <assignment_id?>     List all available assignments if no assignment ID is specified
    submit|s <assignment_id>   Create a new submission to the specified assignment
    help|h                     List all available commands
    version|v                  Displays version info
    
  Example:
  dpc get aed-project          Displays information related to the assignment with 'aed-project' ID
  dpc submit aed-project .     Submits the contents of the current directoryto the assignment with 'aed-project' ID
```

<h2 id="configuration">🔧 Configuration</h2>

```yaml
headless: # it's recommended to store these 2 files in separate locations
  keyFilePath: "/home/{USER}/.dp-cli/"
  keyFileName: "dp-cli.key"
  tokenFilePath: "/home/{USER}/.dp-cli/"
  tokenFileName: "dp-cli.token"
api:
  server: "dp.example.com" # set the server domain here
  endpoints:
    assignments: "/drop-project/api/student/assignments/{ASSIGNMENT_ID}"
    current_assignment: "/drop-project/api/student/assignments/current"
    submission: "/drop-project/api/student/submissions/{SUBMISSION_ID}"
    new_submission: "/drop-project/api/student/submissions/new"
dev: # this is extremelly insecure! should be removed ASAP to force encryption!
  forceFallbackMethod: false
  enableTokenEncryption: true
```
<h2 id="disclaimer">❗ Disclaimer</h2>

The **Drop Project CLI** is not affiliated with, endorsed by, or in any way officially connected with **Drop Project organization**.
The official **Drop Project** website can be found at <a href="https://dropproject.org/">dropproject.org</a>

---

> GitHub [@GmsGarcia](https://github.com/GmsGarcia)
