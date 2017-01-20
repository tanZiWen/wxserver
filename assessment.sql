CREATE TABLE if not EXISTS p_assessment (
  id        BIGINT PRIMARY KEY,
  mobile    varchar(11),
  name varchar(100)
  q1  varchar(10) ,
  q2  varchar(10),
  q3  varchar(10),
  q4  varchar(10),
  q5  varchar(10),
  q6  varchar(10),
  q7  varchar(10),
  q8  varchar(10),
  q9  varchar(10),
  q10 varchar(10),
  q11 varchar(10),
  q12 varchar(10),
  score INTEGER,
  crt TIMESTAMP,
  del boolean
);

COMMENT ON TABLE p_assessment IS '评估表';
COMMENT ON COLUMN p_assessment.id IS '主键';
COMMENT ON COLUMN p_assessment.mobile IS '电话';
COMMENT ON COLUMN p_assessment.q1 IS '问题1';
COMMENT ON COLUMN p_assessment.q2 IS '问题2';
COMMENT ON COLUMN p_assessment.q3 IS '问题3';
COMMENT ON COLUMN p_assessment.q4 IS '问题4';
COMMENT ON COLUMN p_assessment.q5 IS '问题5';
COMMENT ON COLUMN p_assessment.q6 IS '问题6';
COMMENT ON COLUMN p_assessment.q7 IS '问题7';
COMMENT ON COLUMN p_assessment.q8 IS '问题8';
COMMENT ON COLUMN p_assessment.q9 IS '问题9';
COMMENT ON COLUMN p_assessment.q10 IS '问题10';
COMMENT ON COLUMN p_assessment.q11 IS '问题11';
COMMENT ON COLUMN p_assessment.q12 IS '问题12';
COMMENT ON COLUMN p_assessment.score IS '得分';
COMMENT ON COLUMN p_assessment.crt IS '创建时间';
COMMENT ON COLUMN p_assessment.del IS '是否删除';

