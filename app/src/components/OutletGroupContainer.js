import React from 'react';
import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';

import OutletGroup from './OutletGroup';
import { apiRequest } from '../util'

const styles = {
  outletGroupContainer: {
    marginTop: 64,
  },
}

class OutletGroupContainer extends React.Component {
  constructor(props, context) {
    super(props, context)

    this.state = {
      outletGroups: [],
    };
  }

  componentDidMount() {
    apiRequest('GET', '/status')
      .then(result => this.setState({ outletGroups: result }))
      .catch(err => console.error(err));
  }

  render() {
    const { classes } = this.props;

    return (
      <div className={classes.outletGroupContainer}>
        {this.state.outletGroups.map(group =>
          <OutletGroup
            key={group.id}
            group={group}
          />
        )}
      </div>
    );
  }
}

OutletGroupContainer.propTypes = {
  classes: PropTypes.object.isRequired,
};

export default withStyles(styles)(OutletGroupContainer);
