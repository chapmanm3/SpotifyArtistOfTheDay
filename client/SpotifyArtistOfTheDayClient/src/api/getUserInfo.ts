import axios from "axios"

interface UserInfoObj {
  ID: number,
  CreatedAt: string,
  UpdatedAt: string,
  DeletedAt: string,
  Country: string,
  DisplayName: string,
  Email: string,
  ExplicitContent: boolean,
  Followers: number,
  ImageUrl: string,
  Uri: string,
  AuthInfo: {
    ID: number,
    CreatedAt: string,
    UpdatedAt: string,
    DeletedAt: string,
    UserInfoID: number,
    AccessToken: string,
    TokenType: string,
    Scope: string,
    ExpiresIn: number,
    RefreshToken: string
  }
}

interface UserInfoResponse {
  user_info: UserInfoObj
}

export default function getUserInfo(): Promise<UserInfoResponse> {
  return axios.get('http://localhost:8080/api/userInfo', { withCredentials: true }).then((resp) => resp.data)
}


