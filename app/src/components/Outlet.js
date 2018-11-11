import React from 'react';
import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';
import ListItem from '@material-ui/core/ListItem';
import ListItemSecondaryAction from '@material-ui/core/ListItemSecondaryAction';
import ListItemText from '@material-ui/core/ListItemText';
import IconButton from '@material-ui/core/IconButton';
import Icon from '@material-ui/core/Icon';
import Switch from '@material-ui/core/Switch';
import cyan from '@material-ui/core/colors/cyan';
import grey from '@material-ui/core/colors/grey';

import TimeSwitchDialog from './TimeSwitchDialog';
import { apiRequest, outletEnabled, formatTime } from '../util';

const styles = {
  buttonTimeSwitchOn: {
    color: cyan[500],
  },
  buttonTimeSwitchOff: {
    color: grey[500],
  },
};

class Outlet extends React.Component {
  state = {
    enabled: false,
    timeSwitchDialogOpen: false,
    timeSwitch: {
      from: null,
      to: null,
      enabled: false,
    },
  };

  constructor(props, context) {
    super(props, context)

    this.props.registerOutlet(this);
  }

  componentWillReceiveProps(nextProps) {
    const outlet = nextProps.attributes;

    this.setState({ enabled: outletEnabled(outlet) });
  }

  handleToggle = () => {
    const { groupId, outletId } = this.props;

    apiRequest('POST', '/outlet', { groupId, outletId, action: 'toggle' })
      .then(outlet => this.setState({ enabled: outletEnabled(outlet) }))
      .catch(err => console.error(err));
  }

  handleTimeSwitchDialogOpen = () => {
    this.setState({ timeSwitchDialogOpen: true });
  }

  handleTimeSwitchDialogClose = timeSwitch => {
    this.setState({ timeSwitchDialogOpen: false });
  }

  handleTimeSwitchDialogApply = timeSwitch => {
    this.setState({
      timeSwitchDialogOpen: false,
      timeSwitch: timeSwitch,
    });
  }

  getTimeSwitchInfo() {
    const { timeSwitch } = this.state;

    if (!timeSwitch.enabled) {
      return;
    }

    return `${formatTime(timeSwitch.from)} - ${formatTime(timeSwitch.to)}`
  }

  render() {
    const { classes } = this.props;
    const { identifier } = this.props.attributes;
    const { enabled, timeSwitch, timeSwitchDialogOpen} = this.state;
    const timeSwitchButtonClass = timeSwitch.enabled
      ? classes.buttonTimeSwitchOn
      : classes.buttonTimeSwitchOff;

    return (
      <ListItem>
        <ListItemText primary={identifier} secondary={this.getTimeSwitchInfo()} />
        <ListItemSecondaryAction>
          <IconButton className={timeSwitchButtonClass} onClick={this.handleTimeSwitchDialogOpen}>
            <Icon>av_timer</Icon>
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
