import React from "react";
import { Switch, Route } from "react-router-dom";
import Login from "./components/Login";
import Users from "./components/users/Users";
import Configs from "./components/configs/Configs";
import Services from "./components/services/Services";

export const routes = [
    {
        name: 'Login',
        exact: true,
        url: '/',
        component: Login,
    },
    {
        name: 'Users',
        exact: false,
        url: '/users',
        component: Users,
    },
    {
        name: 'Configs',
        exact: false,
        url: '/configs',
        component: Configs,
    },
    {
        name: 'Services',
        exact: false,
        url: '/services',
        component: Services,
    }
]

class MainRouter extends React.Component {
    render() {
        return (
            <Switch>
                {routes.map((route, i) =>
                    (<Route
                        key={i}
                        exact={route.exact}
                        path={route.url}
                        component={route.component}
                    />)
                )}
            </Switch>
        )
    }
}

export default MainRouter;