import { createBrowserRouter } from "react-router";
import Home from "./pages/Home";
import Layout from "./Layout";

const router = createBrowserRouter([
  {
    path: "/",
    Component: Layout,
    children: [
    	{ path: "/", Component: Home },
    ]
  }
]);

export default router
