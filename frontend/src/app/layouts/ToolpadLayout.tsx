import { Link, Outlet} from "react-router-dom";
import * as React from 'react';
import { styled, useTheme, Theme, CSSObject } from '@mui/material/styles';
import MuiDrawer from '@mui/material/Drawer';
import MuiAppBar, { AppBarProps as MuiAppBarProps } from '@mui/material/AppBar';
import { 
    ListItemText, 
    ListItemIcon, 
    ListItemButton, 
    ListItem, 
    Toolbar, 
    List,
    Typography,
    CssBaseline,
    Divider,
    IconButton,
    Box
} from '@mui/material';
import ChevronLeftIcon from '@mui/icons-material/ChevronLeft';
import ChevronRightIcon from '@mui/icons-material/ChevronRight';
import MenuIcon from '@mui/icons-material/Menu';
import MailIcon from '@mui/icons-material/Mail';
import MapsHomeWorkIcon from '@mui/icons-material/MapsHomeWork';
import AnalyticsIcon from '@mui/icons-material/Analytics';
import LogoutIcon from '@mui/icons-material/Logout';
import HomeIcon from '@mui/icons-material/Home';
import { PATHS } from "../../shared/consts";
import authStore from "../../features/auth/model/store";

interface AppBarProps extends MuiAppBarProps {
  open?: boolean;
}

interface ItemPanel {
    label: string;
    icon: React.ReactNode;
    href: string;
};

interface secItemPanel extends ItemPanel {
  action: () => void;
}

const drawerWidth = 240;

const openedMixin = (theme: Theme): CSSObject => ({
  width: drawerWidth,
  transition: theme.transitions.create('width', {
    easing: theme.transitions.easing.sharp,
    duration: theme.transitions.duration.enteringScreen,
  }),
  overflowX: 'hidden',
});

const closedMixin = (theme: Theme): CSSObject => ({
  transition: theme.transitions.create('width', {
    easing: theme.transitions.easing.sharp,
    duration: theme.transitions.duration.leavingScreen,
  }),
  overflowX: 'hidden',
  width: `calc(${theme.spacing(7)} + 1px)`,
  [theme.breakpoints.up('sm')]: {
    width: `calc(${theme.spacing(8)} + 1px)`,
  },
});

const DrawerHeader = styled('div')(({ theme }) => ({
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'flex-end',
  padding: theme.spacing(0, 1),
  // necessary for content to be below app bar
  ...theme.mixins.toolbar,
}));

const AppBar = styled(MuiAppBar, {
  shouldForwardProp: (prop) => prop !== 'open',
})<AppBarProps>(({ theme }) => ({
  zIndex: theme.zIndex.drawer + 1,
  transition: theme.transitions.create(['width', 'margin'], {
    easing: theme.transitions.easing.sharp,
    duration: theme.transitions.duration.leavingScreen,
  }),
  variants: [
    {
      props: ({ open }) => open,
      style: {
        marginLeft: drawerWidth,
        width: `calc(100% - ${drawerWidth}px)`,
        transition: theme.transitions.create(['width', 'margin'], {
          easing: theme.transitions.easing.sharp,
          duration: theme.transitions.duration.enteringScreen,
        }),
      },
    },
  ],
}));

const Drawer = styled(MuiDrawer, { shouldForwardProp: (prop) => prop !== 'open' })(
  ({ theme }) => ({
    width: drawerWidth,
    flexShrink: 0,
    whiteSpace: 'nowrap',
    boxSizing: 'border-box',
    variants: [
      {
        props: ({ open }) => open,
        style: {
          ...openedMixin(theme),
          '& .MuiDrawer-paper': openedMixin(theme),
        },
      },
      {
        props: ({ open }) => !open,
        style: {
          ...closedMixin(theme),
          '& .MuiDrawer-paper': closedMixin(theme),
        },
      },
    ],
  }),
);

const Menu: ItemPanel[] = [
    {
        label: "Объекты",
        icon: <MapsHomeWorkIcon/>,
        href: PATHS.BUILDINGS
    },
    {
        label: "Заявки",
        icon: <MailIcon/>,
        href: PATHS.APPLICATIONS
    },
        {
        label: "Аналитика",
        icon: <AnalyticsIcon/>,
        href: PATHS.DASHBOARD
    }
];

const secMenu: secItemPanel[] = [
    {
        label: "На главную",
        icon: <HomeIcon/>,
        href: PATHS.MAIN,
        action: () => {}
    },
    {
        label: "Выйти",
        icon: <LogoutIcon/>,
        href: `../${PATHS.LOGIN}`,
        action: () => authStore.logout()
    }
];


export default function ToolpadLayout() {
  const theme = useTheme();
  const [open, setOpen] = React.useState(false);

  const handleDrawerOpen = () => {
    setOpen(true);
  };

  const handleDrawerClose = () => {
    setOpen(false);
  };

  return (
    <Box sx={{ display: 'flex' }}>
      <CssBaseline />
      <AppBar position="fixed" open={open}>
        <Toolbar>
          <IconButton
            color="inherit"
            aria-label="open drawer"
            onClick={handleDrawerOpen}
            edge="start"
            sx={[ {marginRight: 5}, open && { display: 'none' } ]}
          >
            <MenuIcon />
          </IconButton>
          <Typography variant="h6" noWrap component="div">
            Панель инструментов
          </Typography>
        </Toolbar>
      </AppBar>
      <Drawer variant="permanent" open={open}>
        <DrawerHeader>
          <IconButton onClick={handleDrawerClose}>
            {theme.direction === 'rtl' ? <ChevronRightIcon /> : <ChevronLeftIcon />}
          </IconButton>
        </DrawerHeader>
        <Divider />
        <List>
          {Menu.map((item) => (
            <ListItem key={item.label} disablePadding sx={{ display: 'block' }}>
              <ListItemButton 
                component={Link}
                to={item.href}
                sx={[
                  {minHeight: 48, px: 2.5},
                  open ? { justifyContent: 'initial' } : { justifyContent: 'center' },
                ]}
              >
                <ListItemIcon
                  sx={[
                    {minWidth: 0, justifyContent: 'center'},
                    open ? { mr: 3 } : { mr: 'auto' },
                  ]}
                >
                  {item.icon}
                </ListItemIcon>
                <ListItemText
                  primary={item.label}
                  sx={[ open ? { opacity: 1 } : { opacity: 0 } ]}
                />
              </ListItemButton>
            </ListItem>
          ))}
          <Divider />
          {secMenu.map((item) => (
            <ListItem key={item.label} disablePadding sx={{ display: 'block' }}>
              <ListItemButton 
                sx={[{ minHeight: 48,px: 2.5 },
                  open ? { justifyContent: 'initial' } : { justifyContent: 'center' },
                ]}
                component={Link}
                to={item.href}
                onClick={item.action}
              >
                <ListItemIcon
                  sx={[{minWidth: 0,justifyContent: 'center'},
                    open ? { mr: 3 } : { mr: 'auto' },
                  ]}
                >
                  {item.icon}
                </ListItemIcon>
                <ListItemText
                  primary={item.label}
                  sx={[ open ? { opacity: 1 } : { opacity: 0 } ]}
                />
              </ListItemButton>
            </ListItem>            
          ))}
        </List>
        <Divider />
      </Drawer>
      <Box component="main" sx={{ flexGrow: 1, padding: "5px" }}>
        <DrawerHeader />
        <Outlet />
      </Box>
    </Box>
  );
}