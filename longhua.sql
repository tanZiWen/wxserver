DROP TABLE IF EXISTS longhua;
CREATE TABLE longhua (
  id     BIGINT PRIMARY KEY,
  name   VARCHAR(50),
  mobile VARCHAR(50),
  doctor VARCHAR(50),
  time   VARCHAR(50),
  crt    TIMESTAMP WITH TIME ZONE,
  lut    TIMESTAMP WITH TIME ZONE,
  status SMALLINT DEFAULT (1)
);

COMMENT ON TABLE longhua IS '龙华膏方节';
COMMENT ON COLUMN longhua.id IS '主键';
COMMENT ON COLUMN longhua.name IS '客户姓名';
COMMENT ON COLUMN longhua.mobile IS '客户手机号码';
COMMENT ON COLUMN longhua.doctor IS '医生姓名';
COMMENT ON COLUMN longhua.time IS '预约时间';
COMMENT ON COLUMN longhua.crt IS '创建时间';
COMMENT ON COLUMN longhua.lut IS '最后更新时间';
COMMENT ON COLUMN longhua.status IS '状态，0:删除, 1:正常';