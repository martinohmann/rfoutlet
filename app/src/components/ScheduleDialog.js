import React from 'react';
import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';
import Dialog from '@material-ui/core/Dialog';
import Button from '@material-ui/core/Button';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import ListItemText from '@material-ui/core/ListItemText';
import Divider from '@material-ui/core/Divider';
import AddIcon from '@material-ui/icons/Add';

import DialogAppBar from './DialogAppBar';
import IntervalListItem from './IntervalListItem';
import IntervalDialog from './IntervalDialog';
import { apiRequest, formatTime, weekdaysShort, scheduleToApp, intervalToApi } from '../util';

const styles = theme => ({
  container: {
    marginTop: 64,
  },
  fab: {
    position: 'absolute',
    bottom: theme.spacing.unit * 2,
    right: theme.spacing.unit * 2,
  },
});

class ScheduleDialog extends React.Component {
  state = {
    open: false,
    editIntervalOpen: false,
    schedule: [],
    currentInterval: null,
  }

  componentDidMount() {
    const { schedule } = this.props;

    this.setState({ schedule: scheduleToApp(schedule) });
  }

  componentWillReceiveProps(nextProps) {
    const { open, schedule } = nextProps;

    this.setState({ open, schedule: scheduleToApp(schedule) });
  }

  handleEditIntervalDialogOpen = (open, currentInterval) => () => {
    this.setState({ editIntervalOpen: open, currentInterval })
  }

  handleIntervalCreate = interval => {
    const { outletId } = this.props;

    this.apiRequest('PUT', { id: outletId, interval: intervalToApi(interval) });
  }

  handleIntervalUpdate = interval => {
    const { outletId } = this.props;

    this.apiRequest('POST', { id: outletId, interval: intervalToApi(interval) });
  }

  handleIntervalDelete = interval => () => {
    const { id } = interval;
    const { outletId } = this.props;

    this.apiRequest('DELETE', { id: outletId, intervalId: id });
  }

  handleIntervalToggle = interval => () => {
    interval.enabled = !interval.enabled;

    this.handleIntervalUpdate(interval);
  }

  apiRequest = (method, data) => {
    apiRequest(method, '/outlet/schedule', data)
      .then(response => response.schedule)
      .then(schedule => this.props.onChange(schedule))
      .catch(err => console.error(err));
  }

  render() {
    const { classes, outletId, onClose } = this.props;
    const { open, editIntervalOpen, schedule, currentInterval } = this.state;

    return (
      <Dialog fullScreen open={open} onClose={onClose}>
        <DialogAppBar
          title="Schedule"
          onClose={onClose}
        />
        <List component="nav" className={classes.container}>
          {schedule.map((interval, i) => {
            return (
              <div key={i}>
                <IntervalListItem
                  interval={interval}
                  onToggle={this.handleIntervalToggle(interval)}
                  onEdit={this.handleEditIntervalDialogOpen(true, interval)}
                  onDelete={this.handleIntervalDelete(interval)}
                />
                <Divider />
              </div>
            )
          })}
          {schedule.length === 0 ? (
            <ListItem>
              <ListItemText primary="No intervals defined yet" />
            </ListItem>
          ) : ''}
        </List>
        <Button variant="fab" color="secondary" className={classes.fab} onClick={this.handleEditIntervalDialogOpen(true, null)}>
          <AddIcon />
        </Button>
        <IntervalDialog
          outletId={outletId}
          open={editIntervalOpen}
          onClose={this.handleEditIntervalDialogOpen(false, null)}
          onSave={this.props.onChange}
          onIntervalCreate={this.handleIntervalCreate}
          onIntervalUpdate={this.handleIntervalUpdate}
          {...currentInterval}
        />
      </Dialog>
    );
  }

  renderInterval = ({ from, to }) => `${formatTime(from)} - ${formatTime(to)}`;

  renderIntervalWeekdays = ({ weekdays }) => weekdays.map(i => weekdaysShort[i]).join(', ');
}

ScheduleDialog.propTypes = {
  classes: PropTypes.object.isRequired,
};

export default withStyles(styles)(ScheduleDialog);
