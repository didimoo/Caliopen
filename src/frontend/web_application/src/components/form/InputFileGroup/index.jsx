import React, { Component } from 'react';
import PropTypes from 'prop-types';
import classnames from 'classnames';
import { InputFile, FieldErrors } from '../';
import Button from '../../Button';

import './style.scss';

const File = ({ file, onRemove, __ }) => (
  <div className="m-input-file-group__file">
    <Button
      className="m-input-file-group__file__remove"
      icon="remove"
      value={file}
      onClick={onRemove}
      shape="plain"
    />
    <span className="m-input-file-group__file__name">{file.name}</span>
    <span className="m-input-file-group__file__size">{ __('input-file-group.file.size', { size: file.size / 1000 }) }</span>
  </div>
);

File.propTypes = {
  onRemove: PropTypes.func.isRequired,
  file: PropTypes.shape({}).isRequired,
  __: PropTypes.func.isRequired,
};

class InputFileGroup extends Component {
  static propTypes = {
    onInputChange: PropTypes.func.isRequired,
    errors: PropTypes.shape({}),
    __: PropTypes.func.isRequired,
    // multiple: PropTypes.bool,
    fileTypes: PropTypes.arrayOf(PropTypes.string),
    className: PropTypes.string,
    descr: PropTypes.string,
  };

  static defaultProps = {
    onChange: null,
    errors: {},
    fileTypes: [''],
    // multiple: false,
    className: null,
    descr: null,
  }

  constructor(props) {
    super(props);
    this.state = {
      file: null,
      fieldError: [],
    };

    this.resetForm = this.resetForm.bind(this);
    this.validate = this.validate.bind(this);
    this.handleInputChange = this.handleInputChange.bind(this);
  }

  handleInputChange(ev) {
    const { __, onInputChange, fileTypes } = this.props;
    const files = ev.target.files.length > 0 ? ev.target.files : null;
    const file = this.validate(files[0]) ? files[0] : null;
    const fieldError = !this.validate(files[0]) ? [__('input-file-group.error.no_valid_ext', { fileTypes: fileTypes.map(type => type) })] : [];
    this.setState({
      file,
      fieldError,
    });

    if (file) { return onInputChange(file); }

    return false;
  }

  resetForm() {
    this.setState({
      file: null,
      fieldError: [],
    });

    return this.props.onInputChange(null);
  }

  validate(file) {
    const fileTypes = this.props.fileTypes;
    const ext = file.name ? `.${file.name.split('.').pop()}` : null;

    return ext && fileTypes.includes(ext);
  }

  render() {
    const { __, errors, descr, className, fileTypes } = this.props;
    const { fieldError, file } = this.state;
    const allErrors = errors ? Object.keys(errors).map(key => errors[key]) : null;

    return (
      <div className={classnames('m-input-file-group', className)}>
        {descr && <p>{descr}</p>}

        {errors.length > 0 && <FieldErrors errors={allErrors} /> }

        {file ?
          <File file={file} onRemove={this.resetForm} __={__} />
          :
          <InputFile
            accept={fileTypes}
            onChange={this.handleInputChange}
            errors={fieldError}
            __={__}
            // multiple={multiple}
          />
        }
      </div>
    );
  }
}

export default InputFileGroup;
