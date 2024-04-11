import { Route, Routes as ReactRoutes } from 'react-router-dom';
import ProtectedRoute from './protected-route';

// Recursive function to flatten nested routes



const generateFlattenRoutes = (routes: any[]): any[] => {
  return routes.reduce((flattenedRoutes, { routes: subRoutes, ...rest }) => {
    const subFlattenRoutes = subRoutes ? generateFlattenRoutes(subRoutes) : [];
    return flattenedRoutes.concat({ ...rest }, subFlattenRoutes);
  }, []);
};
export const renderRoutes = (mainRoutes: any[]) => {
  
  const Routes = ({ }) => {
    const layouts = mainRoutes.map(({ layout: Layout, routes }, index) => {
     
      const subRoutes = generateFlattenRoutes(routes);

      return (
        <Route key={index} element={<Layout />}>
          <Route element={<ProtectedRoute/>}>
            {subRoutes.map(({ component: Component, path, name }) => {
              return Component && path && <Route key={name} element={<Component />} path={path} />;
            })}
          </Route>
        </Route>
      );
    });
    return <ReactRoutes>{layouts}</ReactRoutes>;
  };
  return Routes;
};
