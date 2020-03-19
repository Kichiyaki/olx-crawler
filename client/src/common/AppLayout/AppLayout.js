import React, { useState } from 'react';
import { makeStyles } from '@material-ui/core/styles';
import CssBaseline from '@material-ui/core/CssBaseline';

import AppHeader from './AppHeader/AppHeader';
import Drawer from './Drawer/Drawer';

const useStyles = makeStyles(theme => ({
  root: {
    display: 'flex'
  },
  spacer: theme.mixins.toolbar,
  content: ({ noPadding }) => ({
    flexGrow: 1,
    padding: noPadding ? 0 : theme.spacing(3),
    height: '100vh',
    overflow: 'auto'
  })
}));

const PageLayout = ({ children, headerProps, noPadding }) => {
  const [isOpen, setIsOpen] = useState(false);
  const classes = useStyles({ noPadding });

  const handleDrawerOpen = () => setIsOpen(!isOpen);

  return (
    <div className={classes.root}>
      <CssBaseline />
      <AppHeader
        {...headerProps}
        onDrawerOpen={handleDrawerOpen}
        isOpen={isOpen}
      />
      <Drawer onOpen={handleDrawerOpen} isOpen={isOpen} />
      <main className={classes.content} id="main-content">
        <div className={classes.spacer} id="header-spacer" />
        {children}
      </main>
    </div>
  );
};

PageLayout.defaultProps = {
  headerProps: {},
  noPadding: false
};

export default PageLayout;
