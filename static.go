package main

var LOGINPAGE = `{{define "login"}}

<!DOCTYPE html>
<html>
  <head>
    <title>Open Science Framework | Sign In</title>
    <link rel="icon" herf="/favicon.ico" type="image/x-icon" />
  </head>

  <body>
    <div id="header">
      <span><h3>Open Science Framework | fakeCAS</h3></span>
      </br>
    </div>

    {{if .NotRegistered}}
    <div id="message">
      <p>This login email has been registered but not confirmed. Please check your email (and spam folder). <a href="http://localhost:5000/resend/">Click here</a> to resend your confirmation email.</p>
    </div>
    {{end}}

    {{if .NotAuthorized}}
    <div id="message">
      <p>The service you attempted to authenticate to is not authorized to use CAS.</p>
    </div>
    {{end}}

    {{if .NotValid}}
      <p>Invalid request. User does not exist or invalid/expired validation key.</p>
    {{end}}
    {{if .LoginForm}}
    <div id="forms">
      <form id="login" action="{{.CASLogin}}" method="post">
        <section>
          <span>Email:</span>&nbsp;&nbsp;
          <input id="username" name="username" type="text" value="" size="">&nbsp;&nbsp;
          <input id="submit" type="submit" value="Sign In"/>
            {{if .NotExist}}
            &nbsp;&nbsp;<span>User does not exist.</span>
            {{end}}
        </section>
        <section hidden>
          <input id="password", name="password", type="password", value="" size=""></br>
          <input id="persistence" type="checkbox" value="true" checked/>
          <label id="for-persistence">Stay Signed In</label>
        </section>
      </form>
    </div>
    {{end}}
    </br>
    <div id="links">
      <section>
        <a id="forgot-password" href="http://localhost:5000/forgotpassword/">Forgot Your Password?</a></br>
        <a id="institution-login" href="http://localhost:5000/login/?campaign=institution&redirect_url=http://localhost:5000/">Login Through Your Institution</a></br>
        <a id="back-to-osf" href="http://localhost:5000/">Back to OSF</a></br>
        <a id="create-account" href="http://localhost:5000/register/">Create Account</a>
      </section>
    </div>
  </body>
</html>

{{end}}`
