import React from 'react';
import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';
import ListItem from '@material-ui/core/ListItem';
import ListItemSecondaryAction from '@material-ui/core/ListItemSecondaryAction';
import ListItemText from '@material-ui/core/ListItemText';
import Switch from '@material-ui/core/Switch';

import { apiRequest, outletEnabled } from '../util'

const styles = {};

class Outlet extends React.Component {
  constructor(props, context) {
    super(props, context)

    this.handleToggle = this.handleToggle.bind(this)

    this.state = {
      isEnabled: outletEnabled(props.attributes),
    };

    this.props.registerOutlet(this);
  }

  componentWillReceiveProps(nextProps) {
    const outlet = nextProps.attributes;

    this.setState({ isEnabled: outletEnabled(outlet) });
  }

  handleToggle(event, isEnabled) {
    const data = {
      action: 'toggle',
      group_id: this.props.groupId,
      outlet_id: this.props.outletId
    };

    apiRequest('POST', '/outlet', data)
      .then(outlet => this.setState({ isEnabled: outletEnabled(outlet) }))
      .catch(err => console.error(err));
  }

  render() {
    return (
      <ListItem button onClick={this.handleToggle}>
        <ListItemText primary={this.props.attributes.identifier} />
        <ListItemSecondaryAction>
          <Switch onChange={this.handleToggle} checked={this.state.isEnabled}
          />
        </ListItemSecondaryAction>
      </ListItem>
    );
  }
}

Outlet.propTypes = {
  classes: PropTypes.object.isRequired,
};

export default withStyles(styles)(Outlet);
