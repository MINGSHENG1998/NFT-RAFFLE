import "./widget.scss";
import KeyboardArrowUpIcon from "@mui/icons-material/KeyboardArrowUp";
import AddCardOutlinedIcon from '@mui/icons-material/AddCardOutlined';
import PaidOutlinedIcon from '@mui/icons-material/PaidOutlined';
import PersonAddAltOutlinedIcon from '@mui/icons-material/PersonAddAltOutlined';
import DriveEtaOutlinedIcon from '@mui/icons-material/DriveEtaOutlined';
const Widget = ({ type }) => {
  let data;

  const amount = 100;
  const diff = 20;

switch(type){
    case "order":
        data={
            title:"ORDERS",
            isMoney: false,
            link:"View all orders",
            icon: <AddCardOutlinedIcon className="icon" />            
        };
        break;        
    case "transaction":
        data={
            title:"TRANSACTIONS",
            isMoney: true,
            link:"View all transactions",
            icon: <PaidOutlinedIcon className="icon" />            
        };
        break;
    case "user":
        data={
            title:"USERS",
            isMoney: false,
            link:"See all users",
            icon: <PersonAddAltOutlinedIcon className="icon" />            
        };
        break;
    case "driver":
        data={
            title:"DRIVERS",
            isMoney: false,
            link:"See all drivers",
            icon: <DriveEtaOutlinedIcon className="icon" />            
        };
        break;
    default:
        break;
}



  return (
    <div className="widget">
      <div className="left">
        <span className="title">{data.title}</span>
        <span className="counter">{data.isMoney && "RM"} {amount}</span>
        <span className="link">{data.link}</span>
      </div>
      <div className="right">
        <div className="percentage positive">
          <KeyboardArrowUpIcon />
          {diff} %
        </div>
        {data.icon}
      </div>
    </div>
  );
};

export default Widget;
