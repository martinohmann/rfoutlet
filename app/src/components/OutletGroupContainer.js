import React from 'react';
import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';

import OutletGroup from './OutletGroup';
import { apiRequest } from '../util';

const styles = {
  container: {
    marginTop: 64,
  },
}

class OutletGroupContainer extends React.Component {
  state = {
    groups: [],
  }

  componentDidMount() {
    apiRequest('GET', '/status')
      .then(groups => this.setState({ groups }))
      .catch(err => console.error(err));
  }

  render() {
    const { classes } = this.props;
    const { groups } = this.state;

    return (
      <div className={classes.container}>
        {groups.map((group, groupId) =>
          <OutletGroup
            key={groupId}
            groupId={groupId}
            attributes={group}
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
