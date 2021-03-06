import React, { Component } from 'react';
import PropTypes from 'prop-types';
import Button from '../../../Button';
import { TextFieldGroup } from '../../../form';

import './style.scss';

class TagSearch extends Component {
  static propTypes = {
    onSubmit: PropTypes.func.isRequired,
    onChange: PropTypes.func,
    __: PropTypes.func.isRequired,
  };

  static defaultProps = {
    onChange: null,
  };

  state = {
    terms: '',
  };

  handleChange = (ev) => {
    const terms = ev.target.value;
    this.setState({ terms });
    if (this.props.onChange) {
      this.props.onChange({ terms });
    }
  }

  handleSubmit = () => {
    if (this.state.terms.length === 0) {
      return;
    }

    this.setState((prevState) => {
      this.props.onSubmit({ tag: prevState.terms });

      return {
        terms: '',
      };
    });
  }

  render() {
    const { __ } = this.props;

    return (
      <div className="m-tags-search">
        <TextFieldGroup
          id="tags-search"
          name="terms"
          value={this.state.terms}
          className="m-tags-search__input"
          label={__('tags.form.search.label')}
          placeholder={__('tags.form.search.placeholder')}
          onChange={this.handleChange}
          showLabelforSr
        />
        <Button className="m-tags-search__button" icon="search" onClick={this.handleSubmit} />
      </div>
    );
  }
}

export default TagSearch;
