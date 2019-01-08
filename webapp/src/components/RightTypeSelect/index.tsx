import * as React from 'react'
import { Select } from 'antd';
import UiStore from '../../stores/UiStore';
const Option = Select.Option;

const uiStore = UiStore.getInstance( );

interface RightTypeSelectProps {
  value?: any;
  onChange?: any;
}

export default class RightTypeSelect extends React.Component<RightTypeSelectProps> {
  state = {
    value: null
  }

  constructor(props: RightTypeSelectProps) {
    super(props);

    this.state = {
      value: this.state.value
    }
  }

  onChange = (value: any) => {
    this.props.onChange(value)
  }

  render() {
    return <Select onChange={this.onChange} value={this.props.value}>
      <Option value={0}>{uiStore.T('RIGHT_TYPE_STANDARD')}</Option>
      <Option value={1}>{uiStore.T('RIGHT_TYPE_ADMIN')}</Option>
    </Select>
  }
}