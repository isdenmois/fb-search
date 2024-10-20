<div class="circle" role="progressbar" />

<style>
  .circle {
    align-self: center;
    animation: 0.75s linear 0s infinite normal none running rotate;
    border-color: var(--searchSelected) transparent var(--searchSelected) var(--searchSelected);
    border-width: 2px;
    border-style: solid;
    border-image: initial;
    border-radius: 50%;
    height: 32px;
    margin-top: 32px;
    width: 32px;
  }
  @keyframes rotate {
    0% {
      transform: rotate(0);
    }
    100% {
      transform: rotate(360deg);
    }
  }
</style>
