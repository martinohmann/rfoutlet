import React from 'react';
import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';
import Fab from '@material-ui/core/Fab';
import { List, ListItem } from './List';
import ListItemText from '@material-ui/core/ListItemText';
import AddIcon from '@material-ui/icons/Add';

import ConfigurationDialog from './ConfigurationDialog';
import IntervalListItem from './IntervalListItem';
import IntervalDialog from './IntervalDialog';
import { intervalToApi } from '../util';

const styles = theme => ({
  fab: {
    position: 'absolute',
    bottom: theme.spacing(2),
    right: theme.spacing(2),
  },
});

class ScheduleDialog extends React.Component {
  state = {
    intervalDialogOpen: false,
    currentInterval: null,
  }

  handleIntervalDialogOpen = (open, currentInterval) => () => {
    this.setState({ intervalDialogOpen: open, currentInterval })
  }

  handleIntervalCreate = interval => this.dispatchMessage('create', interval);

  handleIntervalUpdate = interval => this.dispatchMessage('update', interval);

  handleIntervalDelete = interval => () => this.dispatchMessage('delete', interval);

  handleIntervalToggle = interval => () => {
    interval.enabled = !interval.enabled;

    this.handleIntervalUpdate(interval);
  }

  dispatchMessage = (action, interval) => {
    const data = {
      action: action,
      id: this.props.outletId,
      interval: intervalToApi(interval),
    }

    this.props.dispatchMessage({ type: 'interval', data });
  }

  render() {
    const { open, classes, onClose, schedule } = this.props;
    const { intervalDialogOpen, currentInterval } = this.state;

    return (
      <ConfigurationDialog title="Schedule" open={open} onClose={onClose}>
        <List>
          {schedule.map((interval, key) => (
            <IntervalListItem
              key={key}
              interval={interval}
              onToggle={this.handleIntervalToggle(interval)}
              onEdit={this.handleIntervalDialogOpen(true, interval)}
              onDelete={this.handleIntervalDelete(interval)}
            />
          ))}
          {schedule.length === 0 ? (
            <ListItem>
              <ListItemText primary="No intervals defined yet" />
            </ListItem>
          ) : ''}
        </List>
        <Fab
          color="secondary"
          className={classes.fab}
          onClick={this.handleIntervalDialogOpen(true, null)}
        >
          <AddIcon />
        </Fab>
        <IntervalDialog
          open={intervalDialogOpen}
          onClose={this.handleIntervalDialogOpen(false, null)}
          onIntervalCreate={this.handleIntervalCreate}
          onIntervalUpdate={this.handleIntervalUpdate}
          {...currentInterval}
        />
      </ConfigurationDialog>
    );
  }
}

ScheduleDialog.propTypes = {
  classes: PropTypes.object.isRequired,
  open: PropTypes.bool.isRequired,
  onClose: PropTypes.func.isRequired,
  schedule: PropTypes.array.isRequired,
  outletId: PropTypes.string.isRequired,
};

export default withStyles(styles)(ScheduleDialog);
