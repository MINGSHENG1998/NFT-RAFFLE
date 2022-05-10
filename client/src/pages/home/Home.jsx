import Navbar from "../../components/navbar/Navbar";
import Sidebar from "../../components/sidebar/Sidebar";
import Widget from "../../components/widget/Widget";
import "./home.scss";

const Home = () => {
  return (
    <div className="home">
      <Sidebar />
      <div className="homeContainer">
        <Navbar />
        <div className="widgets">
          <Widget type="order"/>
          <Widget type="transaction"/>
          <Widget type="user"/>
          <Widget type="driver"/>
        </div>
        <div className="charts"></div>
      </div>
    </div>
  );
};

export default Home;
