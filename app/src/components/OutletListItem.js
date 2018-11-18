import React from 'react';
import ListItem from '@material-ui/core/ListItem';
import ListItemSecondaryAction from '@material-ui/core/ListItemSecondaryAction';
import ListItemText from '@material-ui/core/ListItemText';
import IconButton from '@material-ui/core/IconButton';
import Switch from '@material-ui/core/Switch';
import EditIcon from '@material-ui/icons/Edit';

import ScheduleDialog from './ScheduleDialog';
import { apiRequest } from '../util';

class OutletListItem extends React.Component {
  state = {
    enabled: false,
    enabledIntervals: [],
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
    let { state, schedule } = outlet;

    schedule = schedule === null ? [] : schedule;

    this.setState({
      enabled: 1 === state,
      enabledIntervals: this.enabledIntervals(schedule),
      schedule,
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
    this.setState({ enabledIntervals: this.enabledIntervals(schedule), schedule })
  }

  enabledIntervals = schedule => {
    return schedule.filter(interval => interval.enabled)
  }

  render() {
    const { id, name } = this.props;
    const { enabled, enabledIntervals, schedule, scheduleDialogOpen } = this.state;

    return (
      <ListItem>
        <ListItemText primary={name} secondary={this.renderIntervals()} />
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

  renderIntervals() {
    const { enabledIntervals } = this.state;

    if (enabledIntervals.length === 0) {
      return ''
    }

    if (enabledIntervals.length === 1) {
      return `1 interval scheduled`
    }

    return `${enabledIntervals.length} intervals scheduled`
  }
}

export default OutletListItem;
