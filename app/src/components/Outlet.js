import React from 'react';
import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';
import ListItem from '@material-ui/core/ListItem';
import ListItemSecondaryAction from '@material-ui/core/ListItemSecondaryAction';
import ListItemText from '@material-ui/core/ListItemText';
import IconButton from '@material-ui/core/IconButton';
import Switch from '@material-ui/core/Switch';
import EditIcon from '@material-ui/icons/Edit';

import ScheduleDialog from './ScheduleDialog';
import { apiRequest } from '../util';

const styles = theme => ({});

class Outlet extends React.Component {
  state = {
    enabled: false,
    schedule: [],
    scheduleDialogOpen: false,
  }

  componentDidMount() {
    this.updateState(this.props);
  }

  componentWillReceiveProps(nextProps) {
    this.updateState(nextProps);
  }

  updateState(outlet) {
    const { state, schedule } = outlet;

    this.setState({
      enabled: 1 === state,
      schedule: schedule === null ? [] : schedule,
    });
  }

  handleToggle = () => {
    const { id } = this.props;

    apiRequest('POST', '/outlet', { id, action: 'toggle' })
      .then(outlet => this.updateState(outlet))
      .catch(err => console.error(err));
  }

  handleScheduleDialogOpen = open => () => {
    this.setState({ scheduleDialogOpen: open });
  }

  handleScheduleChange = schedule => {
    this.setState({ schedule })
  }

  enabledIntervals() {
    const { schedule } = this.state;

    return schedule.filter(interval => interval.enabled)
  }

  renderScheduleText() {
    const enabledIntervals = this.enabledIntervals();

    if (enabledIntervals.length === 0) {
      return ''
    }

    if (enabledIntervals.length === 1) {
      return `1 interval scheduled`
    }

    return `${enabledIntervals.length} intervals scheduled`
  }

  render() {
    const { id, name } = this.props;
    const { enabled, schedule, scheduleDialogOpen } = this.state;
    const enabledIntervals = this.enabledIntervals();

    return (
      <ListItem>
        <ListItemText primary={name} secondary={this.renderScheduleText()} />
        <ListItemSecondaryAction>
          <Switch
            color="primary"
            onChange={this.handleToggle}
            checked={enabled}
            disabled={enabledIntervals.length > 0}
          />
          <IconButton onClick={this.handleScheduleDialogOpen(true)}>
            <EditIcon />
          </IconButton>
        </ListItemSecondaryAction>
        <ScheduleDialog
          outletId={id}
          schedule={schedule}
          open={scheduleDialogOpen}
          onClose={this.handleScheduleDialogOpen(false)}
          onChange={this.handleScheduleChange}
        />
      </ListItem>
    );
  }
}

Outlet.propTypes = {
  classes: PropTypes.object.isRequired,
};

export default withStyles(styles)(Outlet);
