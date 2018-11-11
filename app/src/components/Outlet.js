import React from 'react';
import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';
import ListItem from '@material-ui/core/ListItem';
import ListItemSecondaryAction from '@material-ui/core/ListItemSecondaryAction';
import ListItemText from '@material-ui/core/ListItemText';
import IconButton from '@material-ui/core/IconButton';
import Icon from '@material-ui/core/Icon';
import Switch from '@material-ui/core/Switch';
import { DateTime } from 'luxon';

import TimeSwitchDialog from './TimeSwitchDialog';
import { apiRequest, formatTime } from '../util';

const styles = theme => ({
  buttonTimeSwitchOn: {
    color: theme.palette.primary[500],
  },
  buttonTimeSwitchOff: {
    color: theme.palette.grey[500],
  },
});

class Outlet extends React.Component {
  state = {
    enabled: false,
    timeSwitchDialogOpen: false,
    timeSwitch: {
      from: null,
      to: null,
      enabled: false,
    },
  }

  componentDidMount() {
    this.updateState(this.props);
  }

  componentWillReceiveProps(nextProps) {
    this.updateState(nextProps);
  }

  updateState(outlet) {
    this.setState({ enabled: 1 === outlet.state });

    if (undefined === outlet.timeSwitch) {
      return;
    }

    const { timeSwitch } = outlet;

    this.setState({
      timeSwitch: {
        enabled: timeSwitch.enabled,
        from: DateTime.fromISO(timeSwitch.from),
        to: DateTime.fromISO(timeSwitch.to),
      },
    });
  }

  handleToggle = () => {
    const { groupId, outletId } = this.props;

    apiRequest('POST', '/outlet', { groupId, outletId, action: 'toggle' })
      .then(outlet => this.updateState(outlet))
      .catch(err => console.error(err));
  }

  handleTimeSwitchDialogOpen = () => {
    this.setState({ timeSwitchDialogOpen: true });
  }

  handleTimeSwitchDialogClose = () => {
    this.setState({ timeSwitchDialogOpen: false });
  }

  handleTimeSwitchDialogApply = timeSwitch => {
    this.setState({
      timeSwitchDialogOpen: false,
      timeSwitch: timeSwitch,
    });
  }

  renderTimeSwitchText() {
    const { timeSwitch } = this.state;

    if (!timeSwitch.enabled) {
      return;
    }

    return `${formatTime(timeSwitch.from)} - ${formatTime(timeSwitch.to)}`
  }

  render() {
    const { classes, identifier } = this.props;
    const { enabled, timeSwitch, timeSwitchDialogOpen } = this.state;
    const timeSwitchButtonClass = timeSwitch.enabled
      ? classes.buttonTimeSwitchOn
      : classes.buttonTimeSwitchOff;

    return (
      <ListItem>
        <ListItemText primary={identifier} secondary={this.renderTimeSwitchText()} />
        <ListItemSecondaryAction>
          <IconButton className={timeSwitchButtonClass} onClick={this.handleTimeSwitchDialogOpen}>
            <Icon>schedule</Icon>
          </IconButton>
          <Switch
            color="primary"
            onChange={this.handleToggle}
            checked={enabled}
            disabled={timeSwitch.enabled}
          />
        </ListItemSecondaryAction>
        <TimeSwitchDialog
          identifier={identifier}
          open={timeSwitchDialogOpen}
          from={timeSwitch.from}
          to={timeSwitch.to}
          onApply={this.handleTimeSwitchDialogApply}
          onClose={this.handleTimeSwitchDialogClose}
        />
      </ListItem>
    );
  }
}

Outlet.propTypes = {
  classes: PropTypes.object.isRequired,
};

export default withStyles(styles)(Outlet);
