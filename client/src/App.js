import Home from "./pages/home/Home";
import Login from "./pages/login/Login";
import { BrowserRouter, Routes, Route } from "react-router-dom";
import Admin_List from "./pages/admin/admin_list/Admin_List";
import Admin_Details from "./pages/admin/admin_details/Admin_Details";
import New_Admin from "./pages/admin/new_admin/New_Admin";
import Driver_List from "./pages/driver/driver_list/Driver_List";
import Driver_Details from "./pages/driver/driver_details/Driver_Details";
import Cust_List from "./pages/customer/cust_list/Cust_List";
import Cust_Details from "./pages/customer/cust_details/Cust_Details";

function App() {
  return (
    <div className="App">
      <BrowserRouter>
        <Routes>
          <Route path="/">
            <Route index element={<Home />} />
            <Route path="login" element={<Login />} />
            <Route path="admin">
              <Route index element={<Admin_List />} />
              <Route path=":adminId" element={<Admin_Details />} />
              <Route path="new" element={<New_Admin />} />
            </Route>
            <Route path="driver">
              <Route index element={<Driver_List />} />
              <Route path=":driverId" element={<Driver_Details />} />
            </Route>
            <Route path="cust">
              <Route index element={<Cust_List />} />
              <Route path=":custId" element={<Cust_Details />} />
            </Route>
          </Route>
        </Routes>
      </BrowserRouter>
    </div>
  );
}

export default App;
