import React from 'react';
import PropTypes from 'prop-types';
import MaterialDialog from '@material-ui/core/Dialog';
import DialogAppBar from './DialogAppBar';

export default function Dialog(props) {
  const {
    children,
    onClose,
    onDone,
    doneButtonDisabled,
    doneButtonText,
    title
  } = props;

  return (
    <MaterialDialog fullScreen open onClose={onClose}>
      <DialogAppBar
        title={title}
        onClose={onClose}
        onDone={onDone}
        doneButtonDisabled={doneButtonDisabled}
        doneButtonText={doneButtonText}
      />
      {children}
    </MaterialDialog>
  );
}

Dialog.propTypes = {
  onClose: PropTypes.func.isRequired,
  onDone: PropTypes.func,
  title: PropTypes.string.isRequired,
  doneButtonDisabled: PropTypes.bool,
  doneButtonText: PropTypes.any,
};
