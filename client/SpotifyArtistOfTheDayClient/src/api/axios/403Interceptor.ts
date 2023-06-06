import axios from 'axios'

export default function FourOhThreeInterceptor() {
  axios.interceptors.response.use(function(resp) {
    return resp;
  }, function(err) {
    window.localStorage.setItem("is_authed", "false")
    return err
  })
}
