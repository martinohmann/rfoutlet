import React from 'react';
import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import ListItemText from '@material-ui/core/ListItemText';
import IconButton from '@material-ui/core/IconButton';
import Icon from '@material-ui/core/Icon';

import Outlet from './Outlet';
import { apiRequest } from '../util';

const styles = theme => ({
  group: {
    paddingTop: 1,
    paddingBottom: 1,
    paddingRight: 6,
    background: theme.palette.grey[100],
  },
  groupIdentifier: {
    flexGrow: 1,
    fontWeight: 700,
    color: theme.palette.grey[800],
  },
  buttonGroupOn: {
    color: theme.palette.primary[700],
  },
  buttonGroupOff: {
    color: theme.palette.secondary.light,
  },
});

class OutletGroup extends React.Component {
  state = {
    outlets: [],
  }

  componentDidMount() {
    const { outlets } = this.props;

    this.setState({ outlets });
  }

  handleButtonClick = action => event => {
    const { groupId } = this.props;

    apiRequest('POST', '/outlet_group', { groupId, action })
      .then(result => result.outlets)
      .then(outlets => this.setState({ outlets }))
      .catch(err => console.error(err));
  }

  render() {
    const { classes, groupId, identifier } = this.props;
    const { outlets } = this.state;

    return (
      <List component="nav">
        <ListItem className={classes.group}>
          <ListItemText primary={identifier} disableTypography={true} className={classes.groupIdentifier} />
          <IconButton className={classes.buttonGroupOff} onClick={this.handleButtonClick('off')}>
            <Icon>power_off</Icon>
          </IconButton>
          <IconButton className={classes.buttonGroupOn} onClick={this.handleButtonClick('on')}>
            <Icon>power</Icon>
          </IconButton>
          <IconButton onClick={this.handleButtonClick('toggle')}>
            <Icon>swap_horiz</Icon>
          </IconButton>
        </ListItem>
        {outlets.map((outlet, outletId) =>
          <Outlet
            key={outletId}
            outletId={outletId}
            groupId={groupId}
            {...outlet}
          />
        )}
      </List>
    );
  }
}

OutletGroup.propTypes = {
  classes: PropTypes.object.isRequired,
};

export default withStyles(styles)(OutletGroup);
