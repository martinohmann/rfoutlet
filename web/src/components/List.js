import React from 'react';
import { makeStyles } from '@material-ui/core/styles';
import MaterialList from '@material-ui/core/List';
import MaterialListItem from '@material-ui/core/ListItem';
import Divider from '@material-ui/core/Divider';

const useStyles = makeStyles(theme => ({
  container: {
    marginTop: 64,
    paddingTop: 0,
  },
}));

export function List(props) {
  const classes = useStyles();

  return <MaterialList component="nav" className={classes.container} {...props} />
}

export function ListItem(props) {
  return (
    <React.Fragment>
      <MaterialListItem {...props} />
      <Divider />
    </React.Fragment>
  )
}