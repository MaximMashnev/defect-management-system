import { createBrowserRouter } from "react-router-dom";
import { PageWrapper } from "../layouts/PageWrapper";
import { RequireAuth } from "../../features/auth/components/RequireAuth";
import { lazy } from "react";
import { PATHS } from "../../shared/consts";

const MainLayout = lazy(() => import("../layouts/MainLayout"));
const ToolpadLayout = lazy(() => import("../layouts/ToolpadLayout"));
const MainPage = lazy(() => import("../../pages/main-page"));
const ApplicationCreatingPage = lazy(() => import("../../pages/application-creating-page"))
const ApplicationPage = lazy(() => import("../../pages/applications-page"));
const BuildingPage = lazy(() => import("../../pages/buildings-page"));
const DashboardPage = lazy(() => import("../../pages/dashboard-page"));
const NotFoundPage = lazy(() => import("../../pages/not-found-page"));
const ApplicationTrackingPage = lazy(() => import("../../pages/application-tracking-page"));
const LoginPage = lazy(() => import("../../pages/login-page"));

export const router = createBrowserRouter([
    {
        path: PATHS.MAIN,
        element: <MainLayout />,
        children: [
            {
                index: true,
                element: <PageWrapper title="Главная страница" component={MainPage} />
            },
            {
                path: PATHS.NEW_APPLICATION,
                element: <PageWrapper title="Создание заявки" component={ApplicationCreatingPage} />
            },
            {
                path: `${PATHS.APPLICATION_TRACKING}/:id`,
                element: <PageWrapper title="Отслеживание заявки" component={ApplicationTrackingPage} />
            }
        ]
    },
    {
        path: PATHS.TOOLPAD,
        element: <RequireAuth component={ToolpadLayout}/>,
        children: [
            {
                path: PATHS.BUILDINGS,
                element: <PageWrapper title="Объекты" component={BuildingPage} />
            },
            {
                path: PATHS.APPLICATIONS,
                element: <PageWrapper title="Заявки" component={ApplicationPage} />
            },
            {
                path: PATHS.DASHBOARD,
                element: <PageWrapper title="Аналитика" component={DashboardPage} />
            },
        ]
    },
    {
        path: PATHS.LOGIN,
        element: <PageWrapper title="Авторизация" component={LoginPage} />
    },
    {
        path: "*",
        element: <PageWrapper title="Страница не найдена" component={NotFoundPage} />
    }
])