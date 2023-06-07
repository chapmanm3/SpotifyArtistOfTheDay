
import './userInfo.css'
import { UserInfoObj } from '../../api/getUserInfo'

interface UserInfoProps {
  user: UserInfoObj
}

const UserInfo = ({user}: UserInfoProps) => {

  return (
    <div className='user-info-container'>

      <div className='user-info-name'>
        {user.DisplayName}
      </div>
      <img
        className='user-info-profile-pic'
        src={`${user.ImageUrl}`}
      />
    </div>
  )
}

export default UserInfo
