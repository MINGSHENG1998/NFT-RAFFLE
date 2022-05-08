import Home from "./pages/home/Home";
import Login from "./pages/login/Login";
import { BrowserRouter, Routes, Route } from "react-router-dom";
import User_List from "./pages/user/user_list/User_List";
import New_Admin from "./pages/user/new_admin/New_Admin";
import Admin_Details from "./pages/user/admin_details/Admin_Details";
import Driver_Details from "./pages/user/driver_details/Driver_Details";
import Cust_Details from "./pages/user/cust_details/Cust_Details";

function App() {
  return (
    <div className="App">
      <BrowserRouter>
        <Routes>
          <Route path="/">
            <Route index element={<Home />} />
            <Route path="login" element={<Login />} />
            <Route path="admin">
              <Route index element={<User_List />} />
              <Route path=":adminId" element={<Admin_Details />} />
              <Route path="new" element={<New_Admin />} />
            </Route>
            <Route path="driver">
              <Route index element={<User_List />} />
              <Route path=":driverId" element={<Driver_Details />} />
            </Route>
            <Route path="cust">
              <Route index element={<User_List />} />
              <Route path=":custId" element={<Cust_Details />} />
            </Route>
          </Route>
        </Routes>
      </BrowserRouter>
    </div>
  );
}

export default App;
