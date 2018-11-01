import React from 'react';

import OutletGroup from './OutletGroup';
import { makeApiRequest } from '../util'
import { styles } from '../config'

class OutletGroupContainer extends React.Component {
  constructor(props, context) {
    super(props, context)

    this.state = {
      outletGroups: [],
    };
  }

  componentDidMount() {
    makeApiRequest('/status', {}, result => {
      this.setState({ outletGroups: result });
    });
  }

  render() {
    return (
      <div style={styles.outletGroupContainer}>
        {this.state.outletGroups.map((attributes, groupId) =>
          <OutletGroup
            key={groupId}
            groupId={groupId}
            attributes={attributes}
          />
        )}
      </div>
    );
  }
}

export default OutletGroupContainer;
