

const Loading = (props) => {
    return (
        <div className="d-flex h-100 justify-content-center align-items-center">
          <div>
            <div style={props.style || {}} className={"loadding " + (props.small ? "loadding-small" : "loadding-norman")}>
            </div>
          </div>
      </div>
    )
}

export default Loading
