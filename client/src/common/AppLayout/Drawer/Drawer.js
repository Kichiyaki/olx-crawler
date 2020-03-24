import React from 'react';
import { bool, func } from 'prop-types';
import classNames from 'classnames';
import { MAIN_PAGE, OBSERVATIONS_PAGE, CONFIG_PAGE } from '@config/routes';

import { makeStyles } from '@material-ui/core/styles';
import { Drawer, List, Divider, IconButton } from '@material-ui/core';
import {
  ChevronLeft,
  Dashboard,
  Notifications as NotificationsIcon,
  Settings as SettingsIcon
} from '@material-ui/icons';
import ListItem from './ListItem';

export const DRAWER_WIDTH = 240;

const useStyles = makeStyles(theme => ({
  toolbarIcon: {
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'flex-end',
    padding: '0 8px',
    ...theme.mixins.toolbar
  },
  drawerPaper: {
    position: 'relative',
    whiteSpace: 'nowrap',
    width: DRAWER_WIDTH,
    transition: theme.transitions.create('width', {
      easing: theme.transitions.easing.sharp,
      duration: theme.transitions.duration.enteringScreen
    })
  },
  drawerPaperClose: {
    overflowX: 'hidden',
    transition: theme.transitions.create('width', {
      easing: theme.transitions.easing.sharp,
      duration: theme.transitions.duration.leavingScreen
    }),
    width: theme.spacing(7),
    [theme.breakpoints.up('sm')]: {
      width: theme.spacing(9)
    }
  }
}));

const MyDrawer = ({ isOpen, onOpen, t }) => {
  const classes = useStyles();
  return (
    <Drawer
      variant="permanent"
      classes={{
        paper: classNames(
          classes.drawerPaper,
          !isOpen && classes.drawerPaperClose
        )
      }}
      open={isOpen}
    >
      <div className={classes.toolbarIcon}>
        <IconButton onClick={onOpen}>
          <ChevronLeft />
        </IconButton>
      </div>
      <Divider />
      <List>
        <ListItem to={MAIN_PAGE} text={t('appLayout.drawer.links.mainPage')}>
          <Dashboard />
        </ListItem>
        <ListItem
          to={OBSERVATIONS_PAGE}
          text={t('appLayout.drawer.links.observationsPage')}
        >
          <NotificationsIcon />
        </ListItem>
        <ListItem
          to={CONFIG_PAGE}
          text={t('appLayout.drawer.links.configPage')}
        >
          <SettingsIcon />
        </ListItem>
      </List>
    </Drawer>
  );
};

MyDrawer.defaultProps = {
  isOpen: false,
  onOpen: () => {}
};

MyDrawer.propTypes = {
  onOpen: func.isRequired,
  isOpen: bool.isRequired,
  t: func.isRequired
};

export default MyDrawer;
