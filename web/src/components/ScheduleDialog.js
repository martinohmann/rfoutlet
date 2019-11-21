import React from 'react';
import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';
import Dialog from '@material-ui/core/Dialog';
import Fab from '@material-ui/core/Fab';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import ListItemText from '@material-ui/core/ListItemText';
import Divider from '@material-ui/core/Divider';
import AddIcon from '@material-ui/icons/Add';

import DialogAppBar from './DialogAppBar';
import IntervalListItem from './IntervalListItem';
import IntervalDialog from './IntervalDialog';
import { intervalToApi } from '../util';

const styles = theme => ({
  container: {
    marginTop: 64,
  },
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
    const { outletId } = this.props;
    const data = { action, id: outletId, interval: intervalToApi(interval) }

    this.props.dispatchMessage({ type: 'interval', data });
  }

  render() {
    const { open, classes, outletId, onClose, schedule } = this.props;
    const { intervalDialogOpen, currentInterval } = this.state;

    return (
      <Dialog fullScreen open={open} onClose={onClose}>
        <DialogAppBar title="Schedule" onClose={onClose} />
        <List component="nav" className={classes.container}>
          {schedule.map((interval, key) => {
            return (
              <div key={key}>
                <IntervalListItem
                  interval={interval}
                  onToggle={this.handleIntervalToggle(interval)}
                  onEdit={this.handleIntervalDialogOpen(true, interval)}
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
        <Fab
          color="secondary"
          className={classes.fab}
          onClick={this.handleIntervalDialogOpen(true, null)}
        >
          <AddIcon />
        </Fab>
        <IntervalDialog
          outletId={outletId}
          open={intervalDialogOpen}
          onClose={this.handleIntervalDialogOpen(false, null)}
          onIntervalCreate={this.handleIntervalCreate}
          onIntervalUpdate={this.handleIntervalUpdate}
          {...currentInterval}
        />
      </Dialog>
    );
  }
}

ScheduleDialog.propTypes = {
  classes: PropTypes.object.isRequired,
};

export default withStyles(styles)(ScheduleDialog);
