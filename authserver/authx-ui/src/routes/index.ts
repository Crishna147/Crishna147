// Layouts
import AccountLayout from "../layouts/AccountLayout";
import ApplicationLayout from "../layouts/ApplicationLayout";

// Pages  --> Account and Application

import Login from "../pages/Account/Login";
import SignUp from "../pages/Account/SignUp";
import ChangePassword from "../pages/Account/ChangePassword";


import Home from "../pages/Application/Home";
import Users from "../pages/Application/Users";
import Roles from "../pages/Application/Roles";
import Devices from "../pages/Application/Devices";
import {renderRoutes} from "./PublicRoutes";


export const routes = [
{
    layout: AccountLayout,
    routes: [
      {
        name: 'login',
        title: 'Login page',
        component: Login,
        path: '/login'
      },
      {
        name: 'sign-up',
        title: 'SignUp page',
        component: SignUp,
        path: '/sign-up'
      },
      {
        name: 'change-password',
        title: 'SignUp page',
        component: ChangePassword,
        path: '/change-password'
      }
    ]
  },
{
    layout: ApplicationLayout,
    routes: [
      {
        name: 'home',
        title: 'Home page',
        component: Home,
        path: '/'
      },
      {
        name: 'users',
        title: 'Users',
        hasSiderLink: true,
        routes: [
          {
            name: 'list-users',
            title: 'List of users',
            hasSiderLink: true,
            component: Users,
            path: '/users'
          }
        //,
        //   {
        //     name: 'create-user',
        //     title: 'Add user',
        //     hasSiderLink: true,
        //     component: CreateUser,
        //     path: '/users/new'
        //   }
        ]
      },
      {
        name: 'roles',
        title: 'Roles',
        hasSiderLink: true,
        routes: [
          {
            name: 'list-roles',
            title: 'List of Roles',
            hasSiderLink: true,
            component: Roles,
            path: '/roles'
          }
        ]
      },
      {
        name: 'devices',
        title: 'Devices',
        hasSiderLink: true,
        routes: [
          {
            name: 'list-devices',
            title: 'List of Devices',
            hasSiderLink: true,
            component: Devices,
            path: '/devices'
          }
        ]
      },
    ]
  }
];

export const Routes = renderRoutes(routes);