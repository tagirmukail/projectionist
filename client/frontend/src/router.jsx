import Login from "./components/Login";
import Users from "./components/users/Users";
import Configs from "./components/configs/Configs";
import Services from "./components/services/Services";
import Logout from "./components/Logout"

export const loginRoute = {
    name: 'Login',
    exact: true,
    url: '/',
    component: Login,
};

export const logoutRoute = {
  name: "Logout",
  exact: true,
  url: '/logout',
  component: Logout,
};

export const routes = [
    {
        name: 'Users',
        exact: true,
        url: '/users',
        component: Users,
    },
    {
        name: 'Configs',
        exact: true,
        url: '/configs',
        component: Configs,
    },
    {
        name: 'Services',
        exact: true,
        url: '/services',
        component: Services,
    },
];