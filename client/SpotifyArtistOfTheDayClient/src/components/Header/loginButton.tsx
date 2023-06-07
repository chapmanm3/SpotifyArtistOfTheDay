import handleLogin from "../../api/handleLogin"
import './loginButton.css'

function LoginButton() {

  return (
    <button
      className="login-button"
      onClick={handleLogin}
    >
      Login
    </button>
  )
}

export default LoginButton
