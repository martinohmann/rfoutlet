import React from 'react';
import { makeStyles } from '@material-ui/core/styles';
import PropTypes from 'prop-types';
import List from '@material-ui/core/List';

import OutletListItem from './OutletListItem';

const useStyles = makeStyles(theme => ({
  container: {
    padding: 0,
  },
}));

export default function OutletList(props) {
  const classes = useStyles();

  return (
    <List className={classes.container}>
      {props.outlets.map(outlet =>
        <OutletListItem key={outlet.id} {...outlet} dispatchMessage={props.dispatchMessage} />
      )}
    </List>
  );
}

OutletList.propTypes = {
 outlets: PropTypes.array.isRequired,
 dispatchMessage: PropTypes.func.isRequired,
};
