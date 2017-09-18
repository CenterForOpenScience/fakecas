package main

var TEMPLATES = `
{{define "login"}}

<!DOCTYPE html>
<html>
  <head>
    <title>Open Science Framework | Sign In</title>
  </head>

  <body>
    <div id="header">
      <span><h3>Open Science Framework | fakeCAS</h3></span>
      <br>
    </div>

    <div id="forms">
      <form id="login" action="{{.CASLogin}}" method="post">
        <section>
          <span>Email:</span>&nbsp;&nbsp;
          <input id="username" name="username" type="text" value="" size="" />&nbsp;&nbsp;
          <input id="submit" type="submit" value="Sign In" />
        </section>
        <section hidden>
          <br>
          <span>Password:</span>&nbsp;&nbsp;
          <input id="password", name="password", type="password", value="" />&nbsp;&nbsp;
          <input id="persistence" name="persistence", type="checkbox" value="true" checked />
          <label id="for-persistence">Stay Signed In</label>
        </section>
        <section>
          <br>
          <div id="message">
            {{if .NotExist}}
              <span>User does not exist.</span>
            {{end}}
            {{if .NotRegistered}}
              <span>Login Failed. This login email has been registered but not confirmed.</span>
            {{end}}
          </div>
        </section>
      </form>
      <br>
    </div>

    <br>
    <div id="links">
      <section>
        <a id="back-to-osf" href={{.OSFDomain}}>Back to OSF</a><br>
        <a id="create-account" href={{.CASRegister}}>Create Account</a>
      </section>
    </div>
  </body>
</html>

{{end}}


{{define "register"}}

<!DOCTYPE html>
<html>
  <head>
    <title>Open Science Framework | Sign Up</title>
  </head>

  <body>
    <div id="header">
      <span><h3>Open Science Framework | fakeCAS</h3></span>
      <br>
    </div>

    <div id="forms">
      <form id="register" action="{{.CASRegister}}" method="post">
        <section>
          <span>Fullname:&nbsp;&nbsp;</span>
          <input id="fullname" name="fullname" type="text" value="" size="" /><br><br>
          <span>Email:&nbsp;&nbsp;</span>
          <input id="email" name="email" type="text" value="" size="" /><br><br>
          <span>Password:&nbsp;&nbsp;</span>
          <input id="password" name="password" type="password" value="" size="" /><br><br>
          <input id="submit" type="submit" value="Create Your OSF Account" /><br><br>
        </section>

        <section hidden>
          <input id="persistence" name="persistence" type="checkbox" value="true" checked />
          <label id="persistence">Stay Signed In</label>
        </section>

        <section>
          {{if .RegisterSuccessful}}
            <span>
              &nbsp;&nbsp;A new OSF account has been created. Please <a href={{.CASLogin}}>log in</a> to continue.
            </span>
            <br>
          {{end}}
          {{if .ShowErrorMessages}}
            <span>
              &nbsp;&nbsp;Request to create an OSF account failed! Please take a look at the fakeCAS log.
            </span><br>
          {{end}}
        </section>
      </form>
      <br>
    </div>

    <div id="links">
      <section>
        <a id="back-to-osf" href={{.OSFDomain}}>Back to OSF</a><br>
        <a id="sign-in" href={{.CASLogin}}>Already have and account?</a><br>
      </section>
      <br>
    </div>
  </body>

</html>

{{end}}



{{define "unauthorized"}}

<!DOCTYPE html>
<html>
  <head>
    <title>Open Science Framework | fakeCAS</title>
  </head>

  <body>
    <div id="header">
      <span><h3>Open Science Framework | fakeCAS</h3></span>
      <br>
    </div>

    <div id="message">
      <p>The service you attempted to authenticate to is not authorized to use CAS.</p>
      <br>
    </div>

    <div id="links">
      <section>
        <a id="back-to-osf" href={{.OSFDomain}}>Back to OSF</a><br>
      </section>
    </div>
  </body>
</html>

{{end}}


{{define "invalid"}}

<!DOCTYPE html>
<html>
  <head>
    <title>Open Science Framework | fakeCAS</title>
  </head>

  <body>
    <div id="header">
      <span><h3>Open Science Framework | fakeCAS</h3></span>
      <br>
    </div>

    <div id="message">
      <p>Invalid request! User does not exist or invalid verification key.</p>
      <br>
    </div>

    <div id="links">
      <section>
        <a id="back-to-osf" href={{.OSFDomain}}>Back to OSF</a><br>
        <a id="sign-in" href={{.CASLogin}}>Already have and account?</a><br>
        <a id="create-account" href={{.CASRegister}}>Create Account</a>
      </section>
    </div>
  </body>
</html>

{{end}}

`
