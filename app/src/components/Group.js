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
  container: {
    paddingTop: 1,
    paddingBottom: 1,
    paddingRight: 6,
    background: theme.palette.grey[100],
  },
  groupName: {
    flexGrow: 1,
    fontWeight: 700,
    color: theme.palette.grey[800],
  },
  buttonOn: {
    color: theme.palette.primary[700],
  },
  buttonOff: {
    color: theme.palette.secondary.light,
  },
});

class Group extends React.Component {
  state = {
    outlets: [],
  }

  componentDidMount() {
    const { outlets } = this.props;

    this.setState({ outlets });
  }

  handleButtonClick = action => event => {
    const { id } = this.props;

    apiRequest('POST', '/outlet_group', { id, action })
      .then(result => result.outlets)
      .then(outlets => this.setState({ outlets }))
      .catch(err => console.error(err));
  }

  render() {
    const { classes, name } = this.props;
    const { outlets } = this.state;

    return (
      <List component="nav">
        <ListItem className={classes.container}>
          <ListItemText className={classes.groupName} primary={name} disableTypography={true} />
          <IconButton className={classes.buttonOff} onClick={this.handleButtonClick('off')}>
            <Icon>power_off</Icon>
          </IconButton>
          <IconButton className={classes.buttonOn} onClick={this.handleButtonClick('on')}>
            <Icon>power</Icon>
          </IconButton>
          <IconButton onClick={this.handleButtonClick('toggle')}>
            <Icon>swap_horiz</Icon>
          </IconButton>
        </ListItem>
        {outlets.map(outlet =>
          <Outlet key={outlet.id} {...outlet} />
        )}
      </List>
    );
  }
}

Group.propTypes = {
  classes: PropTypes.object.isRequired,
};

export default withStyles(styles)(Group);
