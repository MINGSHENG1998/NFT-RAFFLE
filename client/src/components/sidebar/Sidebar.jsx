import "./sidebar.scss";
import DashboardOutlinedIcon from '@mui/icons-material/DashboardOutlined';
import AssessmentOutlinedIcon from '@mui/icons-material/AssessmentOutlined';
import CreditCardOutlinedIcon from '@mui/icons-material/CreditCardOutlined';
import DriveEtaOutlinedIcon from '@mui/icons-material/DriveEtaOutlined';
import PersonOutlineOutlinedIcon from '@mui/icons-material/PersonOutlineOutlined';
import SupervisorAccountOutlinedIcon from '@mui/icons-material/SupervisorAccountOutlined';
import NotificationsOutlinedIcon from '@mui/icons-material/NotificationsOutlined';
import SettingsSystemDaydreamOutlinedIcon from '@mui/icons-material/SettingsSystemDaydreamOutlined';
import ArticleOutlinedIcon from '@mui/icons-material/ArticleOutlined';
import SettingsOutlinedIcon from '@mui/icons-material/SettingsOutlined';
import AccountCircleOutlinedIcon from '@mui/icons-material/AccountCircleOutlined';
import LogoutOutlinedIcon from '@mui/icons-material/LogoutOutlined';

const Sidebar = () => {
  return (
    <div className="sidebar">
      <div className="sidebar_logo">
        <span className="logo">logo</span>
      </div>
      <hr/>
      <div className="sidebar_list">
        <ul>
          <p className="title">Main</p>
          <li>
            <DashboardOutlinedIcon className="icon"/>
            <span>Dashboard</span>
          </li>
          <li>
            <AssessmentOutlinedIcon className="icon"/>
            <span>Report</span>
          </li>
          <li>
            <CreditCardOutlinedIcon className="icon"/>
            <span>Order</span>
          </li>
          <p className="title">Users</p>
          <li>
            <DriveEtaOutlinedIcon className="icon"/>
            <span>Driver</span>
          </li>
          <li>
            <PersonOutlineOutlinedIcon className="icon"/>
            <span>Customer</span>
          </li>
          <li>
            <SupervisorAccountOutlinedIcon className="icon"/>
            <span>Admin</span>
          </li>
          <p className="title">Service</p>
          <li>
            <NotificationsOutlinedIcon className="icon"/>
            <span>Notification</span>
          </li>
          <p className="title">System</p>
          <li>
            <SettingsSystemDaydreamOutlinedIcon className="icon"/>
            <span>System Health</span>
          </li>
          <li>
            <ArticleOutlinedIcon className="icon"/>
            <span>Logs</span>
          </li>
          <li>
            <SettingsOutlinedIcon className="icon"/>
            <span>Settings</span>
          </li>
          <p className="title">Account</p>
          <li>
            <AccountCircleOutlinedIcon className="icon"/>
            <span>Profile</span>
          </li>
          <li>
            <LogoutOutlinedIcon className="icon"/>
            <span>Logout</span>
          </li>
        </ul>
      </div>
      <div className="sidebar_theme_option">
        <div className="theme_option"></div>
        <div className="theme_option"></div>
      </div>
    </div>
  );
};

export default Sidebar;
